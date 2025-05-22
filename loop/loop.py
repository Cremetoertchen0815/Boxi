import os
import time
from multiprocessing import Process
from PIL import Image
import st7735

# Load images
def load_images(folder, width, height):
    images = []
    filenames = sorted([
        f for f in os.listdir(folder)
        if f.endswith('.png') and f[:-4].isdigit()
    ], key=lambda f: int(f[:-4]))

    for file in filenames:
        path = os.path.join(folder, file)
        try:
            img = Image.open(path).resize((width, height)).convert('RGB')
            images.append(img)
        except Exception as e:
            print(f"Failed to load {file}: {e}")
    return images

def display_loop(images, port, dc_pin, name, core_id, barrier=None):
    # Set CPU affinity
    try:
        os.sched_setaffinity(0, {core_id})
        print(f"{name}: Affinity set to CPU {core_id}")
    except Exception as e:
        print(f"{name}: Failed to set CPU affinity: {e}")

    display = st7735.ST7735(
        port=port,
        cs=st7735.BG_SPI_CS_BACK,
        dc=dc_pin,
        rotation=90,
        invert=False,
        spi_speed_hz=30000000
    )
    display.begin()

    frame_rate = 25
    frame_delay = 1.0 / frame_rate

    while True:
        loop_start = time.perf_counter()

        for img in images:
            start = time.perf_counter()
            display.display(img)
            elapsed = time.perf_counter() - start
            sleep = frame_delay - elapsed
            if sleep > 0:
                time.sleep(sleep)

        loop_elapsed = time.perf_counter() - loop_start
        fps = len(images) / loop_elapsed
        print(f"{name}: Loop FPS = {fps:.2f}")

        if barrier:
            barrier.wait()

# Main setup
if __name__ == "__main__":
    WIDTH, HEIGHT = 160, 128
    frames = load_images("frames", WIDTH, HEIGHT)

    if not frames:
        print("No frames found.")
        exit(1)

    from multiprocessing import Barrier
    sync_barrier = Barrier(2)

    p1 = Process(target=display_loop, args=(frames, 0, "GPIO23", "Display 1", 1, sync_barrier))
    p2 = Process(target=display_loop, args=(frames, 1, "GPIO12", "Display 2", 2, sync_barrier))

    p1.start()
    p2.start()

    p1.join()
    p2.join()
