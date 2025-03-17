def main() -> None:
    path = input('Введите путь к файлу: ')
    target = input('Введите слово для поиска: ')
    res = read_file(path, target)
    print(f'Всего слов в файле: {res[0]}\nСлово {target} повторяется в тексте {res[1]} раз')

def read_file(path: str='foo.txt', target: str='coffee') -> tuple[int, int]:
    try:
        k = 0
        with open(path, 'r', encoding='utf-8') as f:
            words = f.read().split()
            for w in words:
                if target.lower() in w.lower() and target:
                    k += 1
        return (len(words), k)
                
    except Exception as e:
        raise FileNotFoundError(f'Cannot read file: {e}')

if __name__ == "__main__":
    main()