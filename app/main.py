import json
from contextlib import asynccontextmanager

from fastapi import FastAPI

from app.api.main import api_router
from app.api.routes import health
from app.clients.http_client import create_http_client
from app.config.settings import settings
from app.core.handlers.handlers import register_exception_handlers
from app.models import Group


@asynccontextmanager
async def lifespan(app: FastAPI):
    with settings.GROUPS_FILE.open("r", encoding="utf-8") as file:
        raw_groups = json.load(file)
    app.state.groups = [Group.model_validate(g) for g in raw_groups]
    async with create_http_client() as client:
        app.state.http = client
        yield


app = FastAPI(
    lifespan=lifespan,
    title=settings.PROJECT_NAME,
    openapi_url=f"{settings.API_V1_STR}/openapi.json",
)

app.include_router(api_router, prefix=settings.API_V1_STR)
app.include_router(health.router)


register_exception_handlers(app)
