from fastapi import APIRouter

from app.api.routes import schedule

api_router = APIRouter()

api_router.include_router(schedule.router)
