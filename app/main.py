import json
from contextlib import asynccontextmanager

from fastapi import FastAPI
from typing_extensions import Dict

from app.api.main import api_router
from app.clients.http_client import create_http_client
from app.config.settings import settings
from app.core.handlers.handlers import register_exception_handlers


@asynccontextmanager
async def lifespan(app: FastAPI):
    with settings.GROUPS_FILE.open("r", encoding="utf-8") as file:
        app.state.groups = json.load(file)
    async with create_http_client() as client:
        app.state.http = client
        yield


app = FastAPI(
    lifespan=lifespan,
    title=settings.PROJECT_NAME,
    openapi_url=f"{settings.API_V1_STR}/openapi.json",
)


@app.get("/")
async def root() -> Dict[str, str]:
    return {"status": "ok"}


app.include_router(api_router, prefix=settings.API_V1_STR)
register_exception_handlers(app)
