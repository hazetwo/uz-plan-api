from typing import List

from bs4 import BeautifulSoup

from app.models import ScheduleEntry


def parse(html: BeautifulSoup) -> List[ScheduleEntry]:
    return html.text
