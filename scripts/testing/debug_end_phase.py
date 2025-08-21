#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
調試END_PHASE邏輯
"""

import sys
import os

# 設置 UTF-8 編碼輸出
if sys.platform.startswith('win'):
    import codecs
    sys.stdout = codecs.getwriter('utf-8')(sys.stdout.buffer, 'strict')
    sys.stderr = codecs.getwriter('utf-8')(sys.stderr.buffer, 'strict')
    os.environ['PYTHONIOENCODING'] = 'utf-8'

import requests
import json

# 設置
BASE_URL = "http://localhost:8004/api/v1"
USER_SERVICE_URL = "http://localhost:8002/api/v1"
BOB_ID = "94b46616-3b46-41b3-81dc-e95f70bfb7d5"

def get_game_status(game_id, token):
    """獲取遊戲狀態"""
    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
    
    # 獲取turn-info
    response = requests.get(f"{BASE_URL}/games/{game_id}/turn-info")
    turn_info = response.json()
    print(f"   🔍 Turn Info 原始響應: {json.dumps(turn_info, indent=2, ensure_ascii=False)}")
    
    # 獲取完整狀態
    full_response = requests.get(f"{BASE_URL}/games/{game_id}", headers=headers)
    full_state = full_response.json()
    
    return turn_info, full_state

def perform_end_phase(game_id, token):
    """執行END_PHASE動作"""
    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
    data = {
        "game_id": game_id,
        "player_id": BOB_ID,
        "action_type": "END_PHASE",
        "action_data": []
    }
    
    response = requests.post(f"{BASE_URL}/games/{game_id}/actions", headers=headers, json=data)
    return response.json()

def main():
    print("🔍 調試 END_PHASE 邏輯")
    
    # 1. 登入
    login_data = {"identifier": "bob", "password": "bobbob"}
    login_response = requests.post(f"{USER_SERVICE_URL}/auth/login", json=login_data)
    token = login_response.json()['data']['access_token']
    
    # 2. 創建遊戲
    with open("C:\\Users\\weilo\\Desktop\\ua\\test_data\\FULL_50_CARDS_DECK.json", "r", encoding="utf-8") as f:
        game_data = json.load(f)
    
    create_response = requests.post(f"{BASE_URL}/games", json=game_data)
    game_id = create_response.json()['data']['game']['id']
    
    # 3. 加入遊戲
    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
    requests.post(f"{BASE_URL}/games/{game_id}/join", headers=headers)
    
    # 4. 調度
    mulligan_data = {"game_id": game_id, "player_id": BOB_ID, "mulligan": False}
    requests.post(f"{BASE_URL}/games/{game_id}/mulligan", headers=headers, json=mulligan_data)
    
    print("\n📊 執行 END_PHASE 前的狀態:")
    turn_info_before, full_state_before = get_game_status(game_id, token)
    
    # 正確解析響應結構
    actual_turn_info_before = turn_info_before['data'] if 'data' in turn_info_before else turn_info_before
    
    print(f"   回合: {actual_turn_info_before.get('turn', 'N/A')}")
    print(f"   階段: {actual_turn_info_before.get('phase', 'N/A')}")
    print(f"   活躍玩家: {'Bob' if actual_turn_info_before.get('is_player1_turn') else 'Kage'}")
    
    # 5. 執行END_PHASE
    print("\n🎯 執行 END_PHASE...")
    end_phase_result = perform_end_phase(game_id, token)
    print(f"   結果: {'成功' if end_phase_result.get('success') else '失敗'}")
    if end_phase_result.get('error'):
        print(f"   錯誤: {end_phase_result['error']}")
    
    print("\n📊 執行 END_PHASE 後的狀態:")
    turn_info_after, full_state_after = get_game_status(game_id, token)
    
    # 正確解析響應結構
    actual_turn_info_after = turn_info_after['data'] if 'data' in turn_info_after else turn_info_after
    
    print(f"   回合: {actual_turn_info_after.get('turn', 'N/A')}")
    print(f"   階段: {actual_turn_info_after.get('phase', 'N/A')}")
    print(f"   活躍玩家: {'Bob' if actual_turn_info_after.get('is_player1_turn') else 'Kage'}")
    
    # 6. 分析
    print("\n📋 分析結果:")
    turn_before = actual_turn_info_before.get('turn', 0)
    phase_before = actual_turn_info_before.get('phase', 'UNKNOWN')
    turn_after = actual_turn_info_after.get('turn', 0)
    phase_after = actual_turn_info_after.get('phase', 'UNKNOWN')
    
    print(f"   變化: turn {turn_before} → {turn_after}, phase {phase_before} → {phase_after}")
    
    if turn_before == 1 and phase_before == "START":
        if turn_after == 1 and phase_after == "MOVE":
            print("   ✅ 正確: START → MOVE (同回合，階段推進)")
        elif turn_after == 2 and phase_after == "START":
            print("   ❌ 錯誤: START → 下一回合的START (錯誤的回合推進)")
        else:
            print(f"   ⚠️  未預期的變化")

if __name__ == "__main__":
    main()