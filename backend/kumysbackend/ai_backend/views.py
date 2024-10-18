from django.shortcuts import render
from rest_framework.views import APIView
from rest_framework.response import Response

class Recommendations(APIView):

    def post(self, request):
        text = request.data["text"]

        return Response({
            "result": "success",
            "recommendations": [],
        })
