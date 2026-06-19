import glob
import json

all_commits = []

for filename in sorted(glob.glob("git_commits*.json")):
    with open(filename, "r", encoding="utf-8") as f:
        for line in f:
            line = line.strip()
            if not line:
                continue

            try:
                all_commits.append(json.loads(line))
            except json.JSONDecodeError as e:
                print(f"Skipping invalid JSON in {filename}: {e}")

with open("all_commits.json", "w", encoding="utf-8") as f:
    json.dump(all_commits, f, indent=2)

print(f"Wrote {len(all_commits)} commits to all_commits.json")
