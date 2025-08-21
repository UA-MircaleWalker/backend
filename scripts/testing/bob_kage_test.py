#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Union Arena 遊戲測試腳本 - Bob vs Kage
實現 docs/testing/BOB_KAGE_GAME_TEST.md 中的完整測試流程
每一步都顯示詳細的遊戲信息，包括手牌數量、deck數量等
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
import time
from typing import Dict, Any, Optional
from dataclasses import dataclass


@dataclass
class GameStats:
    """遊戲統計信息"""
    player1_hand_count: int = 0
    player2_hand_count: int = 0
    player1_deck_count: int = 0
    player2_deck_count: int = 0
    player1_life_area_count: int = 0
    player2_life_area_count: int = 0
    current_turn: int = 0
    current_phase: str = ""
    active_player: str = ""
    game_status: str = ""


class UAGameTester:
    """Union Arena 遊戲測試器"""
    
    def __init__(self):
        self.base_url = "http://localhost:8004/api/v1"
        self.user_service_url = "http://localhost:8002/api/v1"
        self.bob_token = None
        self.kage_token = None
        self.game_id = None
        self.bob_id = "94b46616-3b46-41b3-81dc-e95f70bfb7d5"
        self.kage_id = "a8e16546-5a86-415a-9baa-ae62b13891b4"
        
    def print_separator(self, title: str):
        """打印分隔線和標題"""
        print(f"\n{'='*60}")
        print(f"  {title}")
        print(f"{'='*60}")
        
    def print_game_stats(self, game_state: Dict[str, Any], title: str = "遊戲狀態"):
        """打印詳細的遊戲統計信息"""
        print(f"\n{title}")
        print("-" * 50)
        
        # 基本遊戲信息
        if 'game' in game_state:
            game = game_state['game']
            print(f"遊戲ID: {game.get('id', 'N/A')}")
            print(f"遊戲狀態: {game.get('status', 'N/A')}")
            print(f"當前回合: {game.get('current_turn', 'N/A')}")
            print(f"當前階段: {game.get('phase', 'N/A')}")
            print(f"活躍玩家: {'Bob' if game.get('active_player') == self.bob_id else 'Kage' if game.get('active_player') == self.kage_id else 'Unknown'}")
            
        # 詳細遊戲狀態
        if 'game_state' in game_state and game_state['game_state']:
            gs = game_state['game_state']
            
            # 新的玩家數據結構 - players 是以ID為鍵的字典
            players = gs.get('players', {})
            
            # 玩家1 (Bob) 統計
            player1 = players.get(self.bob_id, {})
            print(f"\nBob (Player1):")
            print(f"   手牌數量: {len(player1.get('hand', []))}")
            print(f"   牌庫數量: {len(player1.get('deck', []))}")
            print(f"   當前AP: {player1.get('ap', 0)}")
            print(f"   最大AP: {player1.get('max_ap', 0)}")
            
            # 玩家2 (Kage) 統計
            player2 = players.get(self.kage_id, {})
            print(f"\nKage (Player2):")
            print(f"   手牌數量: {len(player2.get('hand', []))}")
            print(f"   牌庫數量: {len(player2.get('deck', []))}")
            print(f"   當前AP: {player2.get('ap', 0)}")
            print(f"   最大AP: {player2.get('max_ap', 0)}")
            
            # 棋盤狀態 - 每個玩家都有自己的棋盤
            player1_board = player1.get('board', {})
            player2_board = player2.get('board', {})
            print(f"\nBob 棋盤狀態:")
            print(f"   前線: {len(player1_board.get('front_line', []))} 張卡")
            print(f"   能源線: {len(player1_board.get('energy_line', []))} 張卡")
            print(f"   墓地: {len(player1_board.get('graveyard', []))} 張卡")
            print(f"   生命區: {len(player1_board.get('life_area', []))} 張卡")
            print(f"   場外區: {len(player1_board.get('outside_area', []))} 張卡")
            print(f"   移除區: {len(player1_board.get('remove_area', []))} 張卡")
            print(f"   公開區: {len(player1_board.get('public_area', []))} 張卡")
            print(f"   隱藏區: {len(player1_board.get('hidden_area', []))} 張卡")
            print(f"\nKage 棋盤狀態:")
            print(f"   前線: {len(player2_board.get('front_line', []))} 張卡")
            print(f"   能源線: {len(player2_board.get('energy_line', []))} 張卡")
            print(f"   墓地: {len(player2_board.get('graveyard', []))} 張卡")
            print(f"   生命區: {len(player2_board.get('life_area', []))} 張卡")
            print(f"   場外區: {len(player2_board.get('outside_area', []))} 張卡")
            print(f"   移除區: {len(player2_board.get('remove_area', []))} 張卡")
            print(f"   公開區: {len(player2_board.get('public_area', []))} 張卡")
            print(f"   隱藏區: {len(player2_board.get('hidden_area', []))} 張卡")
            
        print("-" * 50)
        
    def make_request(self, method: str, url: str, headers: Dict = None, data: Dict = None) -> Dict[str, Any]:
        """發送HTTP請求並處理響應"""
        try:
            if headers is None:
                headers = {"Content-Type": "application/json"}
                
            response = requests.request(method, url, headers=headers, json=data, timeout=10)
            
            print(f"請求: {method} {url}")
            print(f"狀態碼: {response.status_code}")
            
            try:
                result = response.json()
                if response.status_code >= 400:
                    print(f"錯誤: {result.get('error', 'Unknown error')}")
                else:
                    print(f"成功")
                return result
            except json.JSONDecodeError:
                print(f"警告: 響應不是JSON格式: {response.text}")
                return {"error": "Invalid JSON response", "status_code": response.status_code}
                
        except requests.exceptions.RequestException as e:
            print(f"請求失敗: {e}")
            return {"error": str(e)}
    
    def login_user(self, username: str, password: str) -> Optional[str]:
        """登入用戶並返回token"""
        url = f"{self.user_service_url}/auth/login"
        data = {"identifier": username, "password": password}
        
        response = self.make_request("POST", url, data=data)
        if 'data' in response and 'access_token' in response['data']:
            token = response['data']['access_token']
            print(f"{username} 登入成功，Token: {token[:50]}...")
            return token
        elif 'access_token' in response:
            print(f"{username} 登入成功，Token: {response['access_token'][:50]}...")
            return response['access_token']
        print(f"{username} 登入失敗")
        return None
    
    def get_auth_headers(self, token: str) -> Dict[str, str]:
        """獲取認證請求頭"""
        return {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {token}"
        }
    
    def create_game(self) -> str:
        """創建遊戲並返回game_id"""
        self.print_separator("步驟1: 創建遊戲")
        
        # 讀取完整50張卡組數據
        with open("C:\\Users\\weilo\\Desktop\\ua\\test_data\\FULL_50_CARDS_DECK.json", "r", encoding="utf-8") as f:
            game_data = json.load(f)
        
        url = f"{self.base_url}/games"
        response = self.make_request("POST", url, data=game_data)
        
        # 檢查響應結構
        if 'data' in response and 'game' in response['data']:
            game_data = response['data']
            self.game_id = game_data['game']['id']
            print(f"遊戲創建成功！Game ID: {self.game_id}")
            self.print_game_stats(game_data, "初始遊戲狀態")
            return self.game_id
        elif 'game' in response:
            self.game_id = response['game']['id']
            print(f"遊戲創建成功！Game ID: {self.game_id}")
            self.print_game_stats(response, "初始遊戲狀態")
            return self.game_id
        return None
    
    def join_game(self, player_name: str, token: str):
        """玩家加入遊戲"""
        url = f"{self.base_url}/games/{self.game_id}/join"
        headers = self.get_auth_headers(token)
        
        response = self.make_request("POST", url, headers=headers)
        print(f"加入遊戲響應: {json.dumps(response, indent=2, ensure_ascii=False)[:300]}...")
        
        # 檢查響應結構並顯示遊戲狀態
        if 'data' in response and 'game' in response['data']:
            print(f"{player_name} 成功加入遊戲")
            self.print_game_stats(response['data'], f"{player_name} 加入後的遊戲狀態")
        elif 'game' in response:
            print(f"{player_name} 成功加入遊戲")
            self.print_game_stats(response, f"{player_name} 加入後的遊戲狀態")
        elif response.get('success'):
            print(f"{player_name} 成功加入遊戲")
            # 如果沒有遊戲狀態，獲取完整遊戲狀態
            full_state_url = f"{self.base_url}/games/{self.game_id}"
            full_state = self.make_request("GET", full_state_url, headers=headers)
            if 'game_state' in full_state:
                self.print_game_stats(full_state, f"{player_name} 加入後的遊戲狀態")
        else:
            print(f"{player_name} 加入失敗")
        
    def check_game_status(self):
        """檢查遊戲狀態"""
        self.print_separator("檢查遊戲狀態")
        
        # 使用公開端點檢查遊戲信息
        url = f"{self.base_url}/game-info/{self.game_id}"
        response = self.make_request("GET", url)
        print(f"遊戲基本信息: {json.dumps(response, indent=2, ensure_ascii=False)}")
        
        # 使用認證端點獲取完整遊戲狀態
        url = f"{self.base_url}/games/{self.game_id}"
        headers = self.get_auth_headers(self.bob_token)
        response = self.make_request("GET", url, headers=headers)
        
        print(f"完整遊戲狀態響應: {json.dumps(response, indent=2, ensure_ascii=False)[:500]}...")
        
        # 檢查響應結構
        if 'data' in response and 'game_state' in response['data']:
            self.print_game_stats(response['data'], "完整遊戲狀態")
        elif 'game_state' in response:
            self.print_game_stats(response, "完整遊戲狀態")
        elif 'data' in response and 'game' in response['data']:
            self.print_game_stats(response['data'], "完整遊戲狀態")
        elif 'game' in response:
            self.print_game_stats(response, "完整遊戲狀態")
        else:
            print("未找到遊戲狀態數據")
            
    def get_turn_info(self):
        """獲取回合信息"""
        url = f"{self.base_url}/games/{self.game_id}/turn-info"
        response = self.make_request("GET", url)
        
        if 'turn' in response:
            print(f"回合信息:")
            print(f"   回合數: {response.get('turn', 'N/A')}")
            print(f"   階段: {response.get('phase', 'N/A')}")
            print(f"   活躍玩家: {'Bob' if response.get('is_player1_turn') else 'Kage' if response.get('is_player2_turn') else 'Unknown'}")
            return response
        return None
    
    def perform_mulligan(self, player_name: str, token: str, mulligan: bool = False):
        """執行調度"""
        player_id = self.bob_id if player_name == "Bob" else self.kage_id
        url = f"{self.base_url}/games/{self.game_id}/mulligan"
        headers = self.get_auth_headers(token)
        data = {
            "game_id": self.game_id,
            "player_id": player_id,
            "mulligan": mulligan
        }
        
        response = self.make_request("POST", url, headers=headers, data=data)
        action = "重抽手牌" if mulligan else "保留手牌"
        print(f"{player_name} {action}")
        
        if 'game_state' in response:
            self.print_game_stats(response, f"{player_name} 調度後的遊戲狀態")
    
    def perform_action(self, action_type: str, player_name: str, token: str, action_data: list = None):
        """執行遊戲動作"""
        if action_data is None:
            action_data = []
            
        player_id = self.bob_id if player_name == "Bob" else self.kage_id
        url = f"{self.base_url}/games/{self.game_id}/actions"
        headers = self.get_auth_headers(token)
        data = {
            "game_id": self.game_id,
            "player_id": player_id,
            "action_type": action_type,
            "action_data": action_data
        }
        
        print(f"\n{player_name} 執行動作: {action_type}")
        response = self.make_request("POST", url, headers=headers, data=data)
        
        if 'game_state' in response and response['game_state']:
            self.print_game_stats({"game_state": response['game_state']}, f"{player_name} 執行 {action_type} 後")
        
        return response
    
    def run_complete_test(self):
        """執行完整測試流程"""
        print("開始 Union Arena 遊戲測試 - Bob vs Kage")
        print("測試基於 docs/testing/BOB_KAGE_GAME_TEST.md")
        
        # 步驟1: 登入獲取tokens
        self.print_separator("步驟1: 用戶登入")
        self.bob_token = self.login_user("bob", "bobbob")
        self.kage_token = self.login_user("kage", "kagekage")
        
        if not self.bob_token or not self.kage_token:
            print("登入失敗，測試中止")
            return
            
        # 步驟2: 創建遊戲
        game_id = self.create_game()
        if not game_id:
            print("遊戲創建失敗，測試中止")
            return
            
        # 步驟3: 玩家加入遊戲
        self.print_separator("步驟2: 玩家加入遊戲")
        self.join_game("Bob", self.bob_token)
        time.sleep(1)
        self.join_game("Kage", self.kage_token)
        
        # 步驟4: 檢查遊戲狀態
        time.sleep(2)
        self.check_game_status()
        
        # 步驟5: 調度階段
        self.print_separator("步驟3: 調度階段 (Mulligan)")
        self.perform_mulligan("Bob", self.bob_token, False)  # Bob保留手牌
        time.sleep(1)
        self.perform_mulligan("Kage", self.kage_token, False)  # Kage保留手牌
        
        # 步驟6: 檢查調度後狀態
        time.sleep(2)
        self.check_game_status()
        
        # 步驟7: 遊戲動作測試
        self.print_separator("步驟4: 遊戲動作測試")
        
        # 獲取當前回合信息
        turn_info = self.get_turn_info()
        
        # 測試各種動作
        self.test_all_actions()
        
        print("\n測試完成！")
        
    def test_all_actions(self):
        """測試所有遊戲動作"""
        
        # 1. 測試抽牌動作
        self.print_separator("動作測試1: 抽牌 (DRAW_CARD)")
        self.perform_action("DRAW_CARD", "Bob", self.bob_token)
        time.sleep(1)
        
        # 2. 測試額外抽牌
        self.print_separator("動作測試2: 額外抽牌 (EXTRA_DRAW)")
        self.perform_action("EXTRA_DRAW", "Bob", self.bob_token)
        time.sleep(1)
        
        # 3. 測試結束階段
        self.print_separator("動作測試3: 結束階段 (END_PHASE)")
        self.perform_action("END_PHASE", "Bob", self.bob_token)
        time.sleep(1)
        
        # 4. 測試出牌動作
        self.print_separator("動作測試4: 出牌 (PLAY_CARD)")
        # 假設出第一張手牌到能源線
        self.perform_action("PLAY_CARD", "Bob", self.bob_token, [0, 1])  # [hand_index, destination]
        time.sleep(1)
        
        # 5. 測試結束回合
        self.print_separator("動作測試5: 結束回合 (END_TURN)")
        self.perform_action("END_TURN", "Bob", self.bob_token)
        time.sleep(1)
        
        # 6. 檢查回合轉換
        self.print_separator("檢查回合轉換")
        self.get_turn_info()
        self.check_game_status()
        
        # 7. Kage的回合測試
        self.print_separator("動作測試6: Kage的回合")
        self.perform_action("DRAW_CARD", "Kage", self.kage_token)
        time.sleep(1)
        self.perform_action("END_PHASE", "Kage", self.kage_token)
        time.sleep(1)
        self.perform_action("END_TURN", "Kage", self.kage_token)
        
        # 8. 錯誤測試 - 非當前回合玩家執行動作
        self.print_separator("錯誤測試: 非當前回合玩家執行動作")
        turn_info = self.get_turn_info()
        if turn_info and turn_info.get('is_player1_turn'):
            # Bob回合時，讓Kage嘗試執行動作
            print("測試: Kage在Bob回合時執行動作 (應該返回403)")
            self.perform_action("DRAW_CARD", "Kage", self.kage_token)
        else:
            # Kage回合時，讓Bob嘗試執行動作
            print("測試: Bob在Kage回合時執行動作 (應該返回403)")
            self.perform_action("DRAW_CARD", "Bob", self.bob_token)
        
        # 9. 無效動作測試
        self.print_separator("錯誤測試: 無效動作類型")
        print("測試: 無效動作類型 (應該返回500)")
        self.perform_action("INVALID_ACTION", "Bob", self.bob_token)
        
        # 10. 測試不存在的遊戲
        self.print_separator("錯誤測試: 不存在的遊戲")
        print("測試: 不存在的遊戲ID (應該返回404)")
        fake_game_id = "00000000-0000-0000-0000-000000000000"
        url = f"{self.base_url}/games/{fake_game_id}/actions"
        headers = self.get_auth_headers(self.bob_token)
        data = {
            "game_id": fake_game_id,
            "player_id": self.bob_id,
            "action_type": "DRAW_CARD",
            "action_data": []
        }
        self.make_request("POST", url, headers=headers, data=data)


def main():
    """主函數"""
    tester = UAGameTester()
    tester.run_complete_test()


if __name__ == "__main__":
    main()