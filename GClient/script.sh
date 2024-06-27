#!/bin/bash

# コンパイルされたGoプログラムの名前
PROGRAM="./GClient"

# 固定パラメータ
DURATION=20

# クライアント数の配列
CLIENT_NUMBERS=(1 2 4 8 16 32 64)

# イテレーションの数
ITERATION_NUMBER=5

# コマンドの配列
COMMANDS=("A" "B" "C")

# 各コマンドに対して処理を実行
for cmd in "${COMMANDS[@]}"; do
    # 出力ファイル
    OUTPUT_FILE="output_${cmd}.txt"

    # ファイルの初期化
    > $OUTPUT_FILE

    echo "Testing Command $cmd" >> $OUTPUT_FILE

    # 各クライアント数で処理
    for clients in "${CLIENT_NUMBERS[@]}"; do
        total_sum=0
        
        echo "Testing with $clients clients" >> $OUTPUT_FILE
        
        # ITERATION_NUMBER回実行してデータを収集
        for ((i=1; i<=ITERATION_NUMBER; i++)); do
            echo "  Run $i of $ITERATION_NUMBER" >> $OUTPUT_FILE
            
            # プログラムを実行し、出力を取得
            temp_output=$(echo -e "$cmd\n$clients\n$DURATION" | $PROGRAM)
            
            # Total number of commands executedの値を抽出
            commands=$(echo "$temp_output" | grep "Total number of commands executed:" | awk '{print $NF}')
            
            total_sum=$((total_sum + commands))
            
            echo "    Commands executed: $commands" >> $OUTPUT_FILE
            sleep 5  # 各実行間に5秒の休止
        done
        
        # 平均を計算し、DURATIONで割る
        average=$(echo "scale=2; $total_sum / $ITERATION_NUMBER" | bc)
        result=$(echo "scale=2; $average / $DURATION" | bc)
        
        # 結果をファイルに出力
        echo "Clients: $clients, Average commands per second: $result" >> $OUTPUT_FILE
        echo "----------------------------------------" >> $OUTPUT_FILE
    done

    echo "Results for Command $cmd have been saved to $OUTPUT_FILE"
    echo "========================================"
done
