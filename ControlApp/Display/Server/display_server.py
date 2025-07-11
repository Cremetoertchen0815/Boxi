import os
import time
import threading
import socket
from time import sleep
import lgpio
from PIL import Image, ImageDraw, ImageFont
import st7735

SERVER_HOST = '192.168.4.1'  # Replace with server IP
SERVER_ID = 0
HEADER_BUFFER_SIZE = 15

ANIMATION_DIR = "animations"
FRAME_RATE = 25
FRAME_DELAY = 1.0 / FRAME_RATE
WIDTH = 160
HEIGHT = 128
LINE_HEIGHT = 12
GPIO_BACKLIGHT_DISABLE = 13
GPIO_DISPLAY_ENABLE = 6
GPIO_DISPLAY_RESET = 26
PWM_FREQ = 100

font_path = "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf"
font = ImageFont.truetype(font_path, 14)
h = lgpio.gpiochip_open(0)

def wrap_text(text_to_display, font_used, max_width):
    lines = []
    if not text_to_display.strip():
        return lines  # Skip empty or whitespace-only text

    words = text_to_display.split()
    line = ""

    for word in words:
        test_line = f"{line} {word}".strip()
        if font_used.getlength(test_line) <= max_width:
            line = test_line
        else:
            if line:
                lines.append(line)
            line = word

    if line:
        lines.append(line)

    return lines


def draw_text_overlay(img, text_to_display):
    if not text_to_display:
        return
    draw = ImageDraw.Draw(img)
    wrapped_lines = wrap_text(text_to_display, font, WIDTH - 10)
    if not wrapped_lines:
        return
    padding = 4
    total_height = len(wrapped_lines) * LINE_HEIGHT + 2 * padding
    y_start = int((HEIGHT * 2 / 3) - total_height / 2)
    draw.rectangle(
        (0, y_start - padding, WIDTH, y_start + total_height - padding),
        fill=(0, 0, 0)
    )
    for i, line in enumerate(wrapped_lines):
        y = y_start + i * LINE_HEIGHT
        draw.text((5, y), line, font=font, fill=(255, 255, 255))


class DisplayWorker:
    def __init__(self, display, name):
        self.display = display
        self.name = name
        self.lock = threading.Lock()
        self.text = ""
        self.anim_name = None

        self.text_dirty = False
        self.anim_dirty = False
        self.stop_flag = False

        self.cached_frame_paths = []
        self.cached_raw_frames = []
        self.cached_live_frames = []
        self.text_dirty_frame = []
        self.thread = threading.Thread(target=self.run, daemon=True)
        self.thread.start()

    def update_text(self, new_text):
        with self.lock:
            if new_text != self.text:
                self.text = new_text
                self.text_dirty = True

    def update_animation(self, new_anim_name):
        with self.lock:
            if new_anim_name != self.anim_name:
                self.anim_name = new_anim_name
                self.anim_dirty = True
                self.text_dirty = True
                self.stop_flag = True

    def run(self):
        while True:
            with self.lock:
                anim_changed = self.anim_dirty
                anim_name = self.anim_name
                self.anim_dirty = False
                self.stop_flag = False

            # Reload animation frames if needed
            if anim_changed and anim_name:
                self.cached_frame_paths = self.load_frame_paths(anim_name)
                self.cached_raw_frames = [None] * len(self.cached_frame_paths)
                self.cached_live_frames = [None] * len(self.cached_frame_paths)

            if not self.cached_live_frames:
                time.sleep(0.1)
                continue

            animation_updated_flag = anim_changed

            while not self.stop_flag:
                for frame_nr in range(len(self.cached_frame_paths)):
                    frame_start = time.perf_counter()

                    with self.lock:
                        if self.stop_flag:
                            break
                        if animation_updated_flag or self.text_dirty:
                            self.text_dirty_frame = [True] * len(self.cached_frame_paths)
                        rerender_text = self.text_dirty_frame[frame_nr]
                        current_text = self.text
                        self.text_dirty = False
                    
                    animation_updated_flag = False

                    live_frame = self.cached_live_frames[frame_nr]
                    if rerender_text or live_frame is None:

                        live_frame = self.cached_raw_frames[frame_nr]
                        if live_frame is None:
                            live_frame = Image.open(self.cached_frame_paths[frame_nr])
                            self.cached_raw_frames[frame_nr] = live_frame
                            self.text_dirty_frame[frame_nr] = False

                        # Draw text on a copy of the base frame only if text changed
                        if rerender_text and current_text != "":
                            live_frame = live_frame.copy()
                            draw_text_overlay(live_frame, current_text)

                        self.cached_live_frames[frame_nr] = live_frame

                    self.display.display(live_frame)

                    elapsed = time.perf_counter() - frame_start
                    time.sleep(max(0.0, FRAME_DELAY - elapsed))

    def load_frame_paths(self, anim_name):
        try:
            a_path = os.path.join(ANIMATION_DIR, anim_name)
            frame_files = sorted(
                [f for f in os.listdir(a_path) if f.endswith('.png') and f[:-4].isdigit()],
                key=lambda f: int(f[:-4])
            )
            return [os.path.join(a_path, f) for f in frame_files]
        except Exception as exx:
            print(f"[{self.name}] Failed to load animation: {exx}")
            return []

class BrightnessManager:
    def __init__(self):
        self._thread = None
        self._lock = threading.Lock()
        self.brightness = 1.0  # Initial brightness value
        lgpio.gpio_claim_output(h, GPIO_BACKLIGHT_DISABLE, 1)

    def _start_new_fade_thread(self, decrement, initial_value):
        stop_event = threading.Event()

        def run():
            self.brightness = initial_value
            interval = 0.002  # 2 ms
            next_time = time.perf_counter()

            while self.brightness > 0 and not stop_event.is_set():
                self.brightness = max(0.0, self.brightness - decrement)
                duty_cycle = int((1 - self.brightness ** 3) * 100)
                lgpio.tx_pwm(h, GPIO_BACKLIGHT_DISABLE, PWM_FREQ, duty_cycle)

                next_time += interval
                sleep_duration = next_time - time.perf_counter()
                if sleep_duration > 0:
                    time.sleep(sleep_duration)
                else:
                    next_time = time.perf_counter()

        with self._lock:
            if self._thread and self._thread.is_alive():
                self._stop_event.set()
                self._thread.join()

            self._stop_event = stop_event
            self._thread = threading.Thread(target=run, daemon=True)
            self._thread.start()

    def start_countdown(self, decrement, initial_value):
        if initial_value == 0 or decrement <= 0:
            self.set_brightness(initial_value)
            return

        if initial_value > 1_000_000 or initial_value < 0:
            raise ValueError("Decrement must be a positive number")

        self._start_new_fade_thread(decrement, initial_value)

    def set_brightness(self, value):
        with self._lock:
            if self._thread and self._thread.is_alive():
                self._stop_event.set()
                self._thread.join()

            self.brightness = max(0.0, min(1.0, value))
            duty_cycle = int((1 - self.brightness ** 3) * 100)
            lgpio.tx_pwm(h, GPIO_BACKLIGHT_DISABLE, PWM_FREQ, duty_cycle)

print("Turn off display backlights...")
lgpio.gpio_claim_output(h, GPIO_DISPLAY_ENABLE, 1)  # Start low
lgpio.gpio_claim_output(h, GPIO_DISPLAY_RESET, 1)  # Start low

brightnessManager = BrightnessManager()
brightnessManager.set_brightness(0)
lgpio.gpio_write(h, GPIO_DISPLAY_ENABLE, 0)
lgpio.gpio_write(h, GPIO_DISPLAY_RESET, 0)
time.sleep(0.250)
lgpio.gpio_write(h, GPIO_DISPLAY_RESET, 1)
time.sleep(0.250)

print("Establishing connection with displays...")
os.nice(-20)  # Requires appropriate privileges (root for -20)

# Setup displays
disp1 = st7735.ST7735(port=0, cs=st7735.BG_SPI_CS_BACK, dc="GPIO24", rotation=90, invert=False, spi_speed_hz=30000000, bgr = False)
disp2 = st7735.ST7735(port=1, cs=st7735.BG_SPI_CS_BACK, dc="GPIO12", rotation=90, invert=False, spi_speed_hz=30000000, bgr = False)
disp1.begin()
disp2.begin()

# Worker threads for each display
worker1 = DisplayWorker(disp1, "Display 1")
worker2 = DisplayWorker(disp2, "Display 2")

worker1.update_animation("0")
worker2.update_animation("1")

sleep(0.25)
brightnessManager.set_brightness(1)
print("Displays connected!")

# Create TCP socket
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

def send_answer(callback_bytes, success):
    if callback_bytes[0] == 0 and callback_bytes[1] == 0 and callback_bytes[2] == 0 and callback_bytes[3] == 0:
        return

    return_bytes = bytes([0xE6, 0x21, callback_bytes[0], callback_bytes[1], callback_bytes[2], callback_bytes[3], 0x01 if success else 0x00])
    sock.send(return_bytes)

try:
    # Connect to the server
    sock.connect((SERVER_HOST, 25621))
    welcome_message = bytes([ord('h'), ord('e'), ord('w'), ord('w'), ord('o'), ord(':'), SERVER_ID])
    sock.send(welcome_message)
    print(f"Connected to {SERVER_HOST}:25621")
except Exception as ex:
    print("Connection error:", ex)
    exit()

print("Display server ready. Waiting for commands...")
# Main send loop
while True:
    try:
        header = sock.recv(HEADER_BUFFER_SIZE)
        if not header or len(header) != 15:
            print("Server disconnected. Closing application...")
            exit()


        if header[0] != ord('y') or header[1] != ord('i') or header[2] != ord('f') or header[3] != ord('f') or header[4] > 5 or header[4] < 1:
            continue

        callback = bytes([header[5], header[6], header[7], header[8]])
        parameter = int.from_bytes([header[9], header[10]], byteorder='big', signed=False)
        payloadLen = int.from_bytes([header[11], header[12], header[13], header[14]], byteorder='big', signed=False)
        payload = sock.recv(payloadLen)
        if not payload or len(payload) != payloadLen:
            continue

        match header[4]:
            case 0x01: #DoesAnimationExist, parameter is expected frameCount
                if payloadLen != 4 or not payload:
                    continue

                animationId = int.from_bytes([payload[0], payload[1], payload[2], payload[3]], byteorder='big', signed=False)
                anim_path = os.path.join(ANIMATION_DIR, str(animationId))
                result = False
                if os.path.isdir(anim_path):
                    files = [f for f in os.listdir(anim_path)]
                    if len(files) == parameter:
                        result = True

                send_answer(callback, result)
            case 0x02: #UploadFrame
                animationId = int.from_bytes([payload[0], payload[1], payload[2], payload[3]], byteorder='big', signed=False)
                payload = payload[4:]

                try:
                    anim_path = os.path.join(ANIMATION_DIR, str(animationId))
                    file_path = f"{anim_path}/{parameter:04d}.png"

                    # Create the directory if it doesn't exist
                    os.makedirs(anim_path, exist_ok=True)

                    # Write the byte array to the file
                    with open(file_path, "wb") as f:
                        f.write(payload)

                    send_answer(callback, True)
                except Exception:
                    send_answer(callback, False)

            case 0x03: #PlayAnimation
                animationId = str(int.from_bytes([payload[0], payload[1], payload[2], payload[3]], byteorder='big', signed=False))

                pp = os.path.join(ANIMATION_DIR, animationId)
                if not os.path.isdir(pp):
                    print("Animation " + animationId + " does not exist.")
                    continue

                if (parameter & 0b01) == 0b01:
                    worker1.update_animation(animationId)
                if (parameter & 0b10) == 0b10:
                    worker2.update_animation(animationId)
            case 0x04: #ShowText
                text = payload.decode('utf-8')
                if (parameter & 0b01) == 0b01:
                    worker1.update_text(text)
                if (parameter & 0b10) == 0b10:
                    worker2.update_text(text)
            case 0x05: #DisplayBrightness
                try:
                    brightness = parameter / float(0xFFFF)
                    rawDecrement = int.from_bytes([payload[0], payload[1]], byteorder='big', signed=False)
                    decrementNumber = rawDecrement / float(0xFFFF)

                    if rawDecrement > 0:
                        brightnessManager.start_countdown(decrementNumber, brightness)
                    else:
                        brightnessManager.set_brightness(brightness)
                except Exception:
                    send_answer(callback, False)


    except Exception as e:
        print("Receive error:", e)
        break
