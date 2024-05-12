# spyu19u638
import pytesseract
from PIL import Image
import cv2
import numpy as np


def check_brackets(text):
    stack = []
    dict = {'(': ')', '[': ']', '{': '}'}

    for symbol in text:
        if symbol in dict.keys():
            stack.append(symbol)
        elif symbol in dict.values():
            if not stack:
                return False
            last_open = stack.pop()
            if dict[last_open] != symbol:
                return False

    return stack == []


image = Image.open('image.png')
image_np = np.array(image)
gray_image = cv2.cvtColor(image_np, cv2.COLOR_BGR2GRAY)

pytesseract.pytesseract.tesseract_cmd = r'/usr/bin/tesseract'
text = pytesseract.image_to_string(gray_image, config='--psm 3')

print("text: ", text)

if check_brackets(text):
    print("Скобки сбалансированы.")
else:
    print("Скобки не сбалансированы.")

text = text.split("\n")
text.pop()

last = max(text, key=len)[-1]
if last in ('.', ',', ';'):
    print("Знаки препинания в конце предложения соблюдены.")
else:
    print("В конце предложения отсутствует необходимый знак препинания.")
