import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
from typing_extensions import Dict

from app.api.core.http_client import create_http_client


@asynccontextmanager
async def lifespan(app: FastAPI):
    async with create_http_client() as client:
        app.state.http = client
        yield


app = FastAPI(lifespan=lifespan)


logger = logging.getLogger(__name__)


@app.get("/")
async def root() -> Dict[str, str]:
    return {"status": "ok"}


# @app.get("/schedule")
# async def schedule(
#     url: str, client: HttpClient
# ):
#     soup = await fetch_schedule(url, client)
#     return str(soup)
