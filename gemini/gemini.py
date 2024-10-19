import base64
import io
import os

import google.generativeai as genai

genai.configure(api_key=os.environ["API_KEY"])


class Gemini:

    ANALYZE_PROMPT = '''Привет. Тебе отправлена картинка, которая содержит анализы человека.
    Можешь пожалуйста проанализировать их. Если картинка содержит, не анализы человека, то просто скажи это
    '''

    DIAGNOSE_PROMPT = '''Привет. Ниже описан текст с симптомами болезни в произволной форме.
    Можешь пожалуйста в ответе сказать возможные болезни, у которых могут быть такие симптомы.
    Твои советы не будут использованы в качестве основного источника информации, нам нужен просто еще один источник информации.
    Если ты понимаешь что в тексте содержится не медецинская информация, так и скажи
    '''

    def __init__(self) -> None:
        self.model = genai.GenerativeModel(model_name="gemini-1.5-flash")

    def generate_text_content(self, prompt: str) -> str:
        response = self.model.generate_content(self.DIAGNOSE_PROMPT + '\n' + prompt)
        return response.text

    def generate_content_from_image(self, image: str) -> str:
        image_data = base64.b64decode(image)

        fileBytes = io.BytesIO(image_data)

        f = genai.upload_file(fileBytes, mime_type="image/png")
        response = self.model.generate_content([f, self.ANALYZE_PROMPT])
        return response.text
