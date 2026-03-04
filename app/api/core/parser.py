from typing import List

from bs4 import BeautifulSoup

from app.models import ScheduleEntry
from app.utils import parse_date, parse_time


def parse_schedule(soup: BeautifulSoup) -> List[ScheduleEntry]:
    table = soup.find(id="table_details")
    if table is None:
        raise

    schedules: List[ScheduleEntry] = []

    for row in table.find_all("tr"):
        cells = [cell.text.strip() for cell in row.find_all("td")]

        if len(cells) < 8:
            continue

        try:
            schedule = ScheduleEntry(
                date=parse_date(cells[0]),
                subGroup=cells[2] or None,
                start=parse_time(cells[3]),
                end=parse_time(cells[4]),
                subject=cells[5],
                type=cells[6],
                teacher=[t.strip() for t in cells[7].split(",")],
                room=cells[8] or None,
            )
        except:
            raise

        schedules.append(schedule)

    return schedules
