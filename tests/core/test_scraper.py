import httpx
import pytest

from app.core.scraper import async_fetch


@pytest.mark.asyncio
async def test_async_fetch(client):
    response = await async_fetch(httpx.URL("https://www.python.org/"), client)
    assert response is not None
