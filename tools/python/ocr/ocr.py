from paddleocr import PaddleOCR, draw_ocr

# Paddleocr目前支持的多语言语种可以通过修改lang参数进行切换
# 例如`ch`, `en`, `fr`, `german`, `korean`, `japan`
ocr = PaddleOCR(use_angle_cls=True, lang="ch")  # need to run only once to download and load model into memory
img_path = "D:/work/128d3d2dda8b1bfa.jpg"
result = ocr.ocr(img_path, cls=True)
for line in result:
    print(line)
