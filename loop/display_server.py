import os
import time
import threading
import json
from PIL import Image, ImageDraw, ImageFont
import st7735

PIPE_PATH = "/tmp/display_pipe"
ANIMATION_DIR = "animations"
FRAME_RATE = 25
FRAME_DELAY = 1.0 / FRAME_RATE

WIDTH = 160
HEIGHT = 128
LINE_HEIGHT = 12

font_path = "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf"
font = ImageFont.truetype(font_path, 14)

def wrap_text(text, font, max_width):
    lines = []
    if not text.strip():
        return lines  # Skip empty or whitespace-only text

    words = text.split()
    line = ""

    for word in words:
        test_line = f"{line} {word}".strip()
        if font.getlength(test_line) <= max_width:
            line = test_line
        else:
            if line:
                lines.append(line)
            line = word

    if line:
        lines.append(line)

    return lines

class DisplayWorker:
    def __init__(self, display, name):
        self.display = display
        self.name = name
        self.lock = threading.Lock()
        self.text = ""
        self.anim_name = None
        self.stop_flag = False
        self.cache = {}
        self.thread = threading.Thread(target=self.run)
        self.thread.daemon = True
        self.thread.start()

    def update(self, anim_name, text):
        with self.lock:
            self.text = text
            self.anim_name = anim_name
            self.stop_flag = True  # Signal to break loop immediately

    def run(self):
        while True:
            current_anim = None
            current_text = ""
            self.stop_flag = False

            with self.lock:
                current_anim = self.anim_name
                current_text = self.text

            if not current_anim:
                time.sleep(0.1)
                continue

            anim_path = os.path.join(ANIMATION_DIR, current_anim)
            frame_files = sorted([
                f for f in os.listdir(anim_path)
                if f.endswith('.png') and f[:-4].isdigit()
            ], key=lambda f: int(f[:-4]))

            cached_frames = []

            # First pass: stream from disk, cache with overlay
            for f in frame_files:
                frame_start = time.perf_counter()

                frame_path = os.path.join(anim_path, f)
                try:
                    img = Image.open(frame_path)

                    if current_text != "":
                        # Draw text ONCE before caching
                        draw = ImageDraw.Draw(img)
                        wrapped_lines = wrap_text(current_text, font, WIDTH - 10)

                        if wrapped_lines:
                            padding = 4
                            total_height = len(wrapped_lines) * LINE_HEIGHT + 2 * padding
                            y_start = int((HEIGHT * 2 / 3) - total_height / 2)

                            # Draw background rectangle
                            draw.rectangle(
                                (0, y_start - padding, WIDTH, y_start + total_height - padding),
                                fill=(0, 0, 0)
                            )

                            for i, line in enumerate(wrapped_lines):
                                y = y_start + i * LINE_HEIGHT
                                draw.text((5, y), line, font=font, fill=(255, 255, 255))

                    self.display.display(img)
                    cached_frames.append(img)

                    elapsed = time.perf_counter() - frame_start
                    sleep_time = FRAME_DELAY - elapsed
                    if sleep_time > 0:
                        time.sleep(sleep_time)

                    if self.stop_flag:
                        break
                except Exception as e:
                    print(f"[{self.name}] Error loading frame: {f} -> {e}")

            # Loop cached frames with accurate timing
            while not self.stop_flag:
                for img in cached_frames:
                    frame_start = time.perf_counter()

                    self.display.display(img)

                    elapsed = time.perf_counter() - frame_start
                    sleep_time = FRAME_DELAY - elapsed
                    if sleep_time > 0:
                        time.sleep(sleep_time)

                    if self.stop_flag:
                        break


# Setup displays
disp1 = st7735.ST7735(port=0, cs=st7735.BG_SPI_CS_BACK, dc="GPIO23", rotation=90, invert=False, spi_speed_hz=50000000)
disp2 = st7735.ST7735(port=1, cs=st7735.BG_SPI_CS_BACK, dc="GPIO12", rotation=90, invert=False, spi_speed_hz=50000000)
disp1.begin()
disp2.begin()

# Worker threads for each display
worker1 = DisplayWorker(disp1, "Display 1")
worker2 = DisplayWorker(disp2, "Display 2")

worker1.update("testcard", "")
worker2.update("doggo", "Boxi DisplayServer V.1")

# Setup named pipe
if not os.path.exists(PIPE_PATH):
    os.mkfifo(PIPE_PATH)

print("Display server ready. Waiting for commands...")

while True:
    with open(PIPE_PATH, "r") as pipe:
        for line in pipe:
            try:
                cmd = json.loads(line.strip())
                screen = cmd.get("screen")
                anim = cmd.get("animation")
                text = cmd.get("text", "")
                print(f"Received command: screen={screen}, anim={anim}, text='{text}'")

                if screen == 1:
                    worker1.update(anim, text)
                elif screen == 2:
                    worker2.update(anim, text)
                else:
                    print("Invalid screen number")

            except json.JSONDecodeError:
                print(f"Invalid JSON: {line.strip()}")
