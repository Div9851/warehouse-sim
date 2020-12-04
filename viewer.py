import sys
import os
import curses


def main(stdscr):
    args = sys.argv
    with open(args[1], "r") as f:
        lines = f.readlines()
    states = []
    state = []
    for line in lines:
        if line == "\n":
            if len(state) > 0:
                states.append(state.copy())
            state.clear()
        else:
            state.append(line)
    if len(state) > 0:
        states.append(state.copy())
    now = 0
    while True:
        stdscr.clear()
        for line in states[now]:
            stdscr.addstr(line)
        c = stdscr.getch()
        if c == curses.KEY_LEFT:
            now = max(now - 1, 0)
        elif c == curses.KEY_RIGHT:
            now = min(now + 1, len(states) - 1)
        elif c == ord('q'):
            break


if __name__ == "__main__":
    curses.wrapper(main)
