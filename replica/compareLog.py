import os
import re

def compare_consensus_values(log1_content, log2_content):
    # 正規表現パターンを使用して、consensusValueの値を抽出する
    pattern = r'consensusValue:\s+(\d+)'

    # ログファイル1からconsensusValueの値を抽出する
    values1 = re.findall(pattern, log1_content)

    # ログファイル2からconsensusValueの値を抽出する
    values2 = re.findall(pattern, log2_content)

    # 抽出したconsensusValueの値が一致しているかどうかを確認する
    return values1 == values2

def check_all_log_files(directory):
    log_files = [file for file in os.listdir(directory) if file.startswith('log') and file.endswith('.txt')]

    if len(log_files) < 2:
        print("比較するログファイルが2つ以上必要です。")
        return

    # 最初のログファイルの内容を取得する
    with open(os.path.join(directory, log_files[0]), "r") as file:
        reference_content = file.read()

    # 残りのログファイルと比較する
    for log_file in log_files[1:]:
        with open(os.path.join(directory, log_file), "r") as file:
            current_content = file.read()

        if not compare_consensus_values(reference_content, current_content):
            print(f"{log_files[0]}と{log_file}のconsensusValueの値が一致していません。")
            return

    print("すべてのログファイルのconsensusValueの値は一致しています。")

# ログファイルのディレクトリを指定する
log_directory = "logs/"

# ログファイルを比較する
check_all_log_files(log_directory)