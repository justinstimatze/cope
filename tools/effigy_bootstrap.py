"""Locate an effigy checkout and put it on sys.path.

card2json.py and run_postproc.py both need the same three things: a path from
$EFFIGY_PATH or a sibling checkout, that path importable, and one error message
naming what to clone when it isn't there. The contract lived in both files and
drifted once already — the sibling-checkout default landed in one copy before
the other. adit flagged the duplicate definition and calque flagged the shared
seam, so it lives here now.

effigy owns the .effigy grammar and is the only thing that parses it:
https://github.com/justinstimatze/effigy
"""

import os
import sys

# A sibling checkout, which is the layout `make rules` assumes.
DEFAULT_EFFIGY_PATH = "../effigy"


def require_effigy(script):
    """Make effigy importable, or exit with instructions.

    script is the tools/ filename to quote in the error, so the suggested
    command line is the one the reader just ran.
    """
    path = os.environ.get("EFFIGY_PATH", DEFAULT_EFFIGY_PATH)
    sys.path.insert(0, path)
    try:
        import effigy.parser  # noqa: F401
    except ImportError:
        sys.exit(
            f"error: no effigy checkout at {path}\n"
            "clone https://github.com/justinstimatze/effigy and point EFFIGY_PATH at it:\n"
            f"  EFFIGY_PATH=/path/to/effigy python3 tools/{script} ..."
        )
    return path
