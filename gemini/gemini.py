import base64
import io
import os

import google.generativeai as genai

genai.configure(api_key=os.environ["API_KEY"])


class Gemini:
    def __init__(self) -> None:
        self.model = genai.GenerativeModel(model_name="gemini-1.5-flash")

    def generate_text_content(self, prompt: str) -> str:
        response = self.model.generate_content(prompt)
        return response.text

    def generate_content_from_image(self, image: str) -> str:
        image_data = base64.b64decode(image)

        fileBytes = io.BytesIO(image_data)

        f = genai.upload_file(fileBytes, mime_type="image/png")
        response = self.model.generate_content([f, "can u analyze for me it?"])
        return response.text
