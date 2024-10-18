from django.shortcuts import render
from rest_framework.views import APIView
from rest_framework.response import Response

from drf_yasg.utils import swagger_auto_schema
from drf_yasg import openapi

import sys
sys.path.append(os.path.join(BASE_DIR, 'gemini'))

from gemini import gemini


class Recommendations(APIView):

    @swagger_auto_schema(
        request_body=openapi.Schema(
            type=openapi.TYPE_OBJECT,
            required=["text"],
            properties={
                "text": openapi.Schema(type=openapi.TYPE_STRING)
            }
        ),
        responses={
            200: openapi.Schema(
                type=openapi.TYPE_OBJECT,
                properties={
                    "result": openapi.Schema(type=openapi.TYPE_STRING),
                    "recommendations":  openapi.Schema(
                        type=openapi.TYPE_ARRAY,
                        items=openapi.Schema(type=openapi.TYPE_STRING),
                    ),
                }
            )
        }
    )
    def post(self, request):
        text = request.data["text"]
        g = gemini.Gemini()
        res = g.generate_text_content(text)

        return Response({
            "result": "success",
            "gemini_result": res,
        })


class Analyze(APIView):

    @swagger_auto_schema(
        request_body=openapi.Schema(
            type=openapi.TYPE_STRING,
        ),
        responses={
            200: openapi.Schema(
                type=openapi.TYPE_OBJECT,
                properties={
                    "result": openapi.Schema(type=openapi.TYPE_STRING),
                    "analytics":  openapi.Schema(
                        type=openapi.TYPE_ARRAY,
                        items=openapi.Schema(type=openapi.TYPE_STRING),
                    ),
                }
            )
        }
    )
    def post(self, request):
        image = request.data

        return Response({
            "result": "success",
            "analytics": [],
        })