import re
import json


def parse(filename):
    with open(filename) as fh:
        sensors = []
        for line in fh:
            if not line:
                continue
            xs = re.findall(r"x=(-?\d+)", line)
            ys = re.findall(r"y=(-?\d+)", line)
            sensor = {
                "Position": {
                    "X": int(xs[0]),
                    "Y": int(ys[0]),
                },
                "Beacon": {
                    "X": int(xs[1]),
                    "Y": int(ys[1]),
                },
            }
            sensors.append(sensor)
    return sensors


if __name__ == "__main__":
    sensors = parse("15.ex")
    with open("15.ex.json", "w") as fh:
        json.dump(sensors, fh, indent=2)

    sensors = parse("15.in")
    with open("15.in.json", "w") as fh:
        json.dump(sensors, fh, indent=None)
