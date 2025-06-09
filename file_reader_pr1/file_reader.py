import os

def main() -> None:
    path = input('Enter the file path: ')
    target = input('Enter the target word: ')
    res = read_file(path, target)
    print(f'Total words in the file: {res[0]}\nThe word {target} appears in the text {res[1]} times')

def read_file(path: str, target: str=None) -> tuple[int, int]:
    if not os.path.exists(path):
        raise FileNotFoundError(f'File not found: {path}')

    try:
        k = 0
        with open(path, 'r', encoding='utf-8') as f:
            words = f.read().split()
            if not target:
                return (len(words), 0)
            for w in words:
                if target.lower() in w.lower():
                    k += 1
        return (len(words), k)
                
    except Exception as e:
        raise RuntimeError(f"Error reading file: {e}")

if __name__ == "__main__":
    main()