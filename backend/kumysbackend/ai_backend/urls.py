from django.urls import path

from . import views

urlpatterns = [
    path("diagnose/", views.Recommendations.as_view(), name="diagnose"),
    path("analyze/", views.Analyze.as_view(), name="analyze"),
]