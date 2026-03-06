from typing import List

from bs4 import BeautifulSoup
from pydantic import ValidationError

from app.core.exceptions import ParsingException
from app.models import Group, ScheduleEntry
from app.utils.date import parse_date, parse_time


def parse_schedule(soup: BeautifulSoup) -> List[ScheduleEntry]:
    table = soup.find(id="table_details")
    if table is None:
        raise ParsingException("Schedule table not found")

    schedule_entries: List[ScheduleEntry] = []

    for row in table.find_all("tr"):
        cells = [cell.text.strip() for cell in row.find_all("td")]
        if len(cells) < 9:
            continue

        try:
            entry = ScheduleEntry(
                date=parse_date(cells[0]),
                # Skipping day
                subGroup=cells[2],
                start=parse_time(cells[3]),
                end=parse_time(cells[4]),
                subject=cells[5],
                type=cells[6],
                # I've seen once multiple teachers added to the same exam
                teacher=[t.strip() for t in (cells[7] or "").split(";")],
                room=cells[8] or None,
            )
        except (ValueError, IndexError, ValidationError) as exc:
            raise ParsingException(f"Invalid schedule row {cells}: {exc}")

        schedule_entries.append(entry)

    return schedule_entries


def parse_groups(soup: BeautifulSoup):
    table = soup.find("table", class_="table table-bordered table-condensed")
    if table is None:
        raise ParsingException("Groups table not found")

    group_entries: List[Group] = []

    for cell in table.find_all("tr"):
        # Only supports IT for now
        name = cell.text.split(" Informatyka /")[0].strip()

        a = cell.find("a", href=True)
        if a is None or "":
            raise ParsingException("Group does not have id")

        id = str(a.get("href")).split("ID=")[-1]

        entry = Group(name=name, id=id)
        group_entries.append(entry)

    return group_entries
