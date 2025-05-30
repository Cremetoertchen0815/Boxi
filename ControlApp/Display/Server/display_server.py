import os
import time
import threading
import socket
from time import sleep

import RPi.GPIO as GPIO
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
GPIO_BACKLIGHT_DISABLE = 29
GPIO_DISPLAY_ENABLE = 31
GPIO_DISPLAY_RESET = 37

font_path = "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf"
font = ImageFont.truetype(font_path, 14)

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

        self.cached_raw_frames = []
        self.last_overlay_text = ""
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
                self.stop_flag = True

    def run(self):
        while True:
            with self.lock:
                anim_changed = self.anim_dirty
                text_changed = self.text_dirty
                anim_name = self.anim_name
                current_text = self.text
                self.anim_dirty = False
                self.text_dirty = False
                self.stop_flag = False

            # Reload animation frames if needed
            if anim_changed and anim_name:
                self.cached_raw_frames = self.load_frames(anim_name)

            if not self.cached_raw_frames:
                time.sleep(0.1)
                continue

            while not self.stop_flag:
                for base_frame in self.cached_raw_frames:
                    frame_start = time.perf_counter()

                    # Draw text on a copy of the base frame only if text changed
                    if text_changed or current_text != self.last_overlay_text:
                        img = base_frame.copy()
                        draw_text_overlay(img, current_text)
                        self.last_overlay_text = current_text
                    else:
                        img = base_frame

                    self.display.display(img)

                    elapsed = time.perf_counter() - frame_start
                    time.sleep(max(0.0, FRAME_DELAY - elapsed))

                    with self.lock:
                        if self.stop_flag:
                            break
                        text_changed = self.text_dirty
                        current_text = self.text
                        self.text_dirty = False

    def load_frames(self, anim_name):
        try:
            a_path = os.path.join(ANIMATION_DIR, anim_name)
            frame_files = sorted(
                [f for f in os.listdir(a_path) if f.endswith('.png') and f[:-4].isdigit()],
                key=lambda f: int(f[:-4])
            )
            return [Image.open(os.path.join(a_path, f)).convert("RGB") for f in frame_files]
        except Exception as exx:
            print(f"[{self.name}] Failed to load animation: {exx}")
            return []


print("Turn off display backlights...")
GPIO.setmode(GPIO.BOARD)
GPIO.setup(GPIO_BACKLIGHT_DISABLE, GPIO.OUT)
GPIO.setup(GPIO_DISPLAY_ENABLE, GPIO.OUT)
GPIO.setup(GPIO_DISPLAY_RESET, GPIO.OUT)



GPIO.output(GPIO_BACKLIGHT_DISABLE, GPIO.HIGH)
GPIO.output(GPIO_DISPLAY_ENABLE, GPIO.LOW) # (active low)
GPIO.output(GPIO_DISPLAY_RESET, GPIO.LOW) # (active low)
time.sleep(0.250)
GPIO.output(GPIO_DISPLAY_RESET, GPIO.HIGH)
time.sleep(0.250)

print("Establishing connection with displays...")
# Setup displays
disp1 = st7735.ST7735(port=0, cs=st7735.BG_SPI_CS_BACK, dc="GPIO23", rotation=90, invert=False, spi_speed_hz=50000000)
disp2 = st7735.ST7735(port=1, cs=st7735.BG_SPI_CS_BACK, dc="GPIO12", rotation=90, invert=False, spi_speed_hz=50000000)
disp1.begin()
disp2.begin()

# Worker threads for each display
worker1 = DisplayWorker(disp1, "Display 1")
worker2 = DisplayWorker(disp2, "Display 2")

worker1.update_animation("testcard")
worker2.update_animation("testcard")

sleep(0.25)
GPIO.output(GPIO_BACKLIGHT_DISABLE, GPIO.LOW)
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


        if header[0] != ord('y') or header[1] != ord('i') or header[2] != ord('f') or header[3] != ord('f') or header[4] > 4 or header[4] < 1:
            continue

        callback = bytes([header[5], header[6], header[7], header[8]])
        parameter = int.from_bytes([header[9], header[10]], byteorder='big', signed=False)
        payloadLen = int.from_bytes([header[11], header[12], header[13], header[14]], byteorder='big', signed=False)
        payload = sock.recv(payloadLen)
        if not payload or len(payload) != payloadLen:
            continue

        match header[4]:
            case 0x01: #DoesAnimationExist, parameter is expected frameCount
                if payloadLen != 4:
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
                animationId = int.from_bytes([payload[0], payload[1], payload[2], payload[3]], byteorder='big', signed=False)
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


    except Exception as e:
        print("Receive error:", e)
        break