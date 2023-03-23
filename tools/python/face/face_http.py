import io

import face_recognition
from flask import Flask, jsonify, request
import requests

app = Flask(__name__)


@app.route('/', methods=['GET'])
def image_has_face():
    url = request.args.get('url')
    print("imgurl==", url)
    file = download_one_pic(url)
    if file is None:
        return jsonify({"code": 1, "error": "请求地址出错"})
    img = face_recognition.load_image_file(io.BytesIO(file))
    face_encodings = face_recognition.face_encodings(img)
    if len(face_encodings) > 0:
        return jsonify({
            "code": 0,
            "found": True,
            "count": len(face_encodings)
        })

    # If no valid image file was uploaded, show the file upload form:
    return jsonify({
        "code": 1,
        "found": False
    })



def download_one_pic(url):
    try:
        response = requests.get(url)
        # response.encoding = 'utf-8'
        if response.status_code == 200:
            # print("header==",header(url))
            # print(response.content)
            return response.content
        return None
    except Exception as e:
        return None


if __name__ == "__main__":
    app.run(host='0.0.0.0', port=5001, debug=True)