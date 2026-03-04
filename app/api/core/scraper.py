import httpx
from bs4 import BeautifulSoup

from app.api.core.exceptions import ScheduleFetchException


async def fetch_schedule(url: str, client: httpx.AsyncClient) -> BeautifulSoup:
    try:
        response = await client.get(url)
        response.raise_for_status()
    except httpx.RequestError as exc:
        raise ScheduleFetchException(f"Request error: {exc}")
    except httpx.HTTPStatusError as exc:
        raise ScheduleFetchException(f"HTTP error: {exc}")

    return BeautifulSoup(response.text, "html.parser")
