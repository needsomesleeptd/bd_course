import cv2
import pytesseract
import sys
import re

def main(image_path):
    # image_path = 'formulas/pos09.png'
    image = cv2.imread(image_path)

    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    _, threshold = cv2.threshold(gray, 150, 255, cv2.THRESH_BINARY)

    custom_config = r'--oem 3 --psm 6'
    text = pytesseract.image_to_string(threshold, lang='eng')
    # text = pytesseract.image_to_string(threshold, config=custom_config)

    # print(text)
    lines = text.split('\n')
    numbered = False
    pattern = re.compile("^[0-9]+\.[0-9]+$")
    for line in lines:
        if '(' in line and ')' in line:
            start_index = line.rfind('(') #line.index('(')
            end_index = line.rfind(')')  # line.index(')')
            number = line[start_index+1:end_index]
            result = re.match(pattern, number)
            # print(result)
            # print("number", number, number.strip())
            if result is not None: #if '.' in number and number.replace('.', '').isdecimal():
                numbered = True
                break

    if numbered:
        print("Ошибок нет")
        return True
    else:
        print("Ошибка: отсутствует нумерация")
        return False


if __name__ == "__main__":
    pytesseract.pytesseract.tesseract_cmd = r'/usr/bin/tesseract'
    main(sys.argv[1])
