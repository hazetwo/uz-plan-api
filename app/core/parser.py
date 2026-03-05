from typing import List

from bs4 import BeautifulSoup

from app.core.exceptions import ParsingException
from app.models import ScheduleEntry
from app.utils import parse_date, parse_time


def parse_schedule(soup: BeautifulSoup) -> List[ScheduleEntry]:
    table = soup.find(id="table_details")
    if table is None:
        raise ParsingException("Schedule table not found")

    schedule_entries: List[ScheduleEntry] = []

    for row in table.find_all("tr"):
        cells = [cell.text.strip() for cell in row.find_all("td")]

        if len(cells) < 8:
            continue

        try:
            entry = ScheduleEntry(
                date=parse_date(cells[0]),
                # Skipping day
                subGroup=cells[2] or None,
                start=parse_time(cells[3]) or None,
                end=parse_time(cells[4]) or None,
                subject=cells[5],
                type=cells[6],
                # I've seen once multiple teachers added to the same exam
                teacher=[t.strip() for t in cells[7].split(";")],
                room=cells[8] or None,
            )
        except (ValueError, IndexError) as exc:
            raise ParsingException(f"Invalid schedule row: {exc}")

        schedule_entries.append(entry)

    return schedule_entries
