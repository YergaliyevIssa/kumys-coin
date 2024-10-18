import os

import google.generativeai as genai

genai.configure(api_key=os.environ["API_KEY"])


class Gemini:
    def __init__(self) -> None:
        self.model = genai.GenerativeModel(model_name="gemini-1.5-flash")

    def generate_text_content(self, prompt: str) -> str:
        response = self.model.generate_content(prompt)
        return response.text
