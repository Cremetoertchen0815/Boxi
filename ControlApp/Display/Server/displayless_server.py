import os
import socket

SERVER_HOST = '127.0.0.1'  # Replace with server IP
SERVER_ID = 1
HEADER_BUFFER_SIZE = 15

ANIMATION_DIR = "animations"
FRAME_RATE = 25
FRAME_DELAY = 1.0 / FRAME_RATE
WIDTH = 160
HEIGHT = 128
LINE_HEIGHT = 12
GPIO_BACKLIGHT_DISABLE = 5

print("Turn off display backlights...")
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
        print("Saas '", payloadLen, "'")
        if payloadLen > 0:
            payload = sock.recv(payloadLen)
            if not payload or len(payload) != payloadLen:
                continue


        print("Sooss '", payloadLen, "'")
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
                    print("Display 1 animation set to: '", animationId, "'")
                if (parameter & 0b10) == 0b10:
                    print("Display 2 animation set to: '", animationId, "'")
            case 0x04: #ShowText
                text = payload.decode('utf-8')
                if (parameter & 0b01) == 0b01:
                    print("Display 1 text set to: '", text, "'")
                if (parameter & 0b10) == 0b10:
                    print("Display 2 text set to: '", text, "'")
            case 0x05: #DisplayBrightness
                try:
                    brightness = int((1 - parameter / float(0xFFFF)) * 1000000)
                    print("Display brightness set to: '", brightness, "'")
                except Exception:
                    send_answer(callback, False)


    except Exception as e:
        print("Receive error:", e)
        break