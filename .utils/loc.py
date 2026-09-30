import os
import sys

printable = bytearray({7, 8, 9, 10, 12, 13, 27} | set(range(0x20, 0x100)) - {0x7F})


def is_binary(filename: str) -> bool:
    with open(filename, "rb") as f:
        b = f.read(1024)
        return bool(b.translate(None, printable))


def count_lines(filename: str) -> tuple[int, int, int]:
    total, blank = 0, 0

    try:
        with open(filename, "r", encoding="utf-8", errors="ignore") as f:
            for line in f:
                total += 1
                if not line.strip():
                    blank += 1

    except (FileNotFoundError, PermissionError) as exc:
        print(f"warning: open {filename}: {exc}")
        return 0, 0, 0

    code = total - blank
    return total, code, blank


def display(root: str, target_dirs: list[str]):
    print(f"{'Path':<40} | {'Total':>6} | {'Code':>6} | {'Blank':>6}")
    print("-" * 67)

    grand_files = 0
    grand_dirs = 0
    grand_total = 0
    grand_code = 0
    grand_blank = 0

    for target in target_dirs:
        target_path = os.path.join(root, target)
        if not os.path.exists(target_path):
            continue

        for current_root, _, files in os.walk(target_path):
            if len(files) == 0:
                continue

            grand_dirs += 1
            dir_total = 0
            dir_code = 0
            dir_blank = 0
            print_lines: list[str] = []

            for filename in files:
                filepath = os.path.join(current_root, filename)

                if not is_binary(filepath) and not filename.endswith("_test.go"):
                    grand_files += 1
                    total, code, blank = count_lines(filepath)

                    grand_total += total
                    grand_code += code
                    grand_blank += blank

                    dir_total += total
                    dir_code += code
                    dir_blank += blank

                    print_lines.append(
                        f"  {filename:<38} | {total:>6} | {code:>6} | {blank:>6}"
                    )

                elif filename.endswith("_test.go"):
                    print_lines.append(f"  {filename:<31} [test] |        |        |")

                else:
                    print_lines.append(f"  {filename:<29} [binary] |        |        |")

            dir_name = os.path.relpath(current_root, root)

            print(f"{dir_name:<40} | {dir_total:>6} | {dir_code:>6} | {dir_blank:>6}")
            for line in print_lines:
                print(line)
            print(f"{' ' * 41}|        |        |")

    print("-" * 67)
    print(f"Total files:        {grand_files}")
    print(f"Total directories:  {grand_dirs}")

    print()

    print(f"Total lines:  {grand_total}")
    print(f"Code lines:   {grand_code} ({round((grand_code / grand_total) * 100)}%)")
    print(f"Blank lines:  {grand_blank} ({round((grand_blank / grand_total) * 100)}%)")


if __name__ == "__main__":
    target_dirs = sys.argv[1:]
    for dir in target_dirs:
        if not os.path.isdir(dir):
            print(f"error: not a directory: {dir}")
            sys.exit(1)

    display(os.getcwd(), target_dirs)
