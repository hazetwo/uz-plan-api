from datetime import date as Date
from datetime import time as Time

from pydantic import BaseModel


# This data structure is completely arbitrary
# as I don't know the exact requirements
# for a schedule to be inserted into the university database.
# It might be the case that an administrator
# can insert it without required values.
class ScheduleEntry(BaseModel):
    date: Date | None = None
    subGroup: str | None = None
    start: Time | None = None
    end: Time | None = None
    subject: str
    type: str
    teacher: list[str]
    room: str | None = None


class Group(BaseModel):
    name: str
    id: str
