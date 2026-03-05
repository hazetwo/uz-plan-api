import json

import httpx
from bs4 import BeautifulSoup

from app.config.settings import settings
from app.core.exceptions import FetchScheduleException, UrlNotFoundException
from app.models import Group


async def async_fetch(
    url: httpx.URL, client: httpx.AsyncClient
) -> BeautifulSoup:
    try:
        response = await client.get(url)
        response.raise_for_status()
    except httpx.RequestError as exc:
        raise FetchScheduleException(f"Request error: {exc}")
    except httpx.HTTPStatusError as exc:
        raise FetchScheduleException(f"HTTP error: {exc}")

    return BeautifulSoup(response.text, "html.parser")


async def fetch_schedule(id: str, client: httpx.AsyncClient) -> BeautifulSoup:
    url = get_url_by_id(id)

    if url is None:
        raise UrlNotFoundException(
            "The provided ID does not map to a valid URL."
        )

    soup = await async_fetch(url, client)

    return soup


def get_url_by_id(id: str) -> httpx.URL | None:
    with settings.GROUPS_FILE.open("r", encoding="utf-8") as file:
        data = json.load(file)

    groups = [Group.model_validate(item) for item in data]

    for group in groups:
        if group.id == id:
            return httpx.URL(
                f"https://plan.uz.zgora.pl/grupy_plan.php?ID={id}"
            )

    return None
