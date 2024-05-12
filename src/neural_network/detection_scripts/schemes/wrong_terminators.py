# zheea18u263
import cv2 as cv
import numpy as np

image = cv.imread("scheme_wrong.png")
gray = cv.cvtColor(image, cv.COLOR_BGR2GRAY)
blur = cv.GaussianBlur(gray, (5, 5), 0)
edges = cv.Canny(blur, 50, 150)
contours, _ = cv.findContours(edges, cv.RETR_TREE, cv.CHAIN_APPROX_NONE)
contour_mask = np.zeros_like(edges)

min_diff = 5000000
for contour in contours:
    x, y, w, h = cv.boundingRect(contour)
    if (w <= 80 or w >= 125) and (h <= 80 or h >= 140): #or len(contour) < 500:
         continue
    cv.rectangle(image, (x, y), (x + w, y + h), (0, 255, 0), 2)
    
    
    cv.drawContours(contour_mask, contour, -1, 255, thickness=cv.FILLED)

    ellipses = cv.fitEllipse(contour)
    center, axes, angle = ellipses
    ellipse_mask = np.zeros_like(gray)
    cv.ellipse(ellipse_mask, ellipses, 255, thickness=cv.FILLED)
    cv.ellipse(ellipse_mask, ellipses, 255, thickness=cv.FILLED)
    cv.ellipse(image, (int(center[0]), int(center[1])), (int(axes[0] / 2), int(axes[1] / 2)), angle, 0, 360, (0, 255, 0), 2)

    difference = cv.absdiff(contour_mask, ellipse_mask)

    difference_count = cv.countNonZero(difference)
    
    if difference_count//h < min_diff:
        min_diff = difference_count//h
    
cv.imwrite('ex3.png', image)
    
if min_diff < 140:
    print("Bad schema!!!!")
else:
    print("Good schema.")
