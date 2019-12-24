package enums

// Daily missions appear to be divided into 11 sets of 3, making for 33 daily missions for each category
// Thus, there are 198 daily missions across all 5 categories.
const (
	// "Take out x enemies!" daily missions
	DailyMissionEnemySet1Pos1 = 1 // 10 enemies
	DailyMissionEnemySet1Pos2 = 2 // 10 enemies
	DailyMissionEnemySet1Pos3 = 3 // 10 enemies

	DailyMissionEnemySet2Pos1 = 4 // 30 enemies
	DailyMissionEnemySet2Pos2 = 5 // 40 enemies
	DailyMissionEnemySet2Pos3 = 6 // 50 enemies

	DailyMissionEnemySet3Pos1 = 7 // 40 enemies
	DailyMissionEnemySet3Pos2 = 8 // 50 enemies
	DailyMissionEnemySet3Pos3 = 9 // 60 enemies

	DailyMissionEnemySet4Pos1 = 10 // 40 enemies
	DailyMissionEnemySet4Pos2 = 11 // 50 enemies
	DailyMissionEnemySet4Pos3 = 12 // 60 enemies

	DailyMissionEnemySet5Pos1 = 13 // 50 enemies
	DailyMissionEnemySet5Pos2 = 14 // 60 enemies
	DailyMissionEnemySet5Pos3 = 15 // 70 enemies

	DailyMissionEnemySet6Pos1 = 16 // 50 enemies
	DailyMissionEnemySet6Pos2 = 17 // 60 enemies
	DailyMissionEnemySet6Pos3 = 18 // 70 enemies

	DailyMissionEnemySet7Pos1 = 19 // 60 enemies
	DailyMissionEnemySet7Pos2 = 20 // 70 enemies
	DailyMissionEnemySet7Pos3 = 21 // 80 enemies

	DailyMissionEnemySet8Pos1 = 22 // 60 enemies
	DailyMissionEnemySet8Pos2 = 23 // 70 enemies
	DailyMissionEnemySet8Pos3 = 24 // 80 enemies

	DailyMissionEnemySet9Pos1 = 25 // 70 enemies
	DailyMissionEnemySet9Pos2 = 26 // 80 enemies
	DailyMissionEnemySet9Pos3 = 27 // 90 enemies

	DailyMissionEnemySet10Pos1 = 28 // 70 enemies
	DailyMissionEnemySet10Pos2 = 29 // 80 enemies
	DailyMissionEnemySet10Pos3 = 30 // 90 enemies

	DailyMissionEnemySet11Pos1 = 31 // 80 enemies
	DailyMissionEnemySet11Pos2 = 32 // 90 enemies
	DailyMissionEnemySet11Pos3 = 33 // 100 enemies

	// "Take out x golden enemies!" daily missions
	DailyMissionGoldEnemySet1Pos1 = 34 // 5 golden enemies
	DailyMissionGoldEnemySet1Pos2 = 35 // 5 golden enemies
	DailyMissionGoldEnemySet1Pos3 = 36 // 5 golden enemies

	DailyMissionGoldEnemySet2Pos1 = 37 // 5 golden enemies
	DailyMissionGoldEnemySet2Pos2 = 38 // 10 golden enemies
	DailyMissionGoldEnemySet2Pos3 = 39 // 15 golden enemies

	DailyMissionGoldEnemySet3Pos1 = 40 // 10 golden enemies
	DailyMissionGoldEnemySet3Pos2 = 41 // 15 golden enemies
	DailyMissionGoldEnemySet3Pos3 = 42 // 20 golden enemies

	DailyMissionGoldEnemySet4Pos1 = 43 // 10 golden enemies
	DailyMissionGoldEnemySet4Pos2 = 44 // 15 golden enemies
	DailyMissionGoldEnemySet4Pos3 = 45 // 20 golden enemies

	DailyMissionGoldEnemySet5Pos1 = 46 // 15 golden enemies
	DailyMissionGoldEnemySet5Pos2 = 47 // 20 golden enemies
	DailyMissionGoldEnemySet5Pos3 = 48 // 25 golden enemies

	DailyMissionGoldEnemySet6Pos1 = 49 // 15 golden enemies
	DailyMissionGoldEnemySet6Pos2 = 50 // 20 golden enemies
	DailyMissionGoldEnemySet6Pos3 = 51 // 25 golden enemies

	DailyMissionGoldEnemySet7Pos1 = 52 // 20 golden enemies
	DailyMissionGoldEnemySet7Pos2 = 53 // 25 golden enemies
	DailyMissionGoldEnemySet7Pos3 = 54 // 30 golden enemies

	DailyMissionGoldEnemySet8Pos1 = 55 // 20 golden enemies
	DailyMissionGoldEnemySet8Pos2 = 56 // 25 golden enemies
	DailyMissionGoldEnemySet8Pos3 = 57 // 30 golden enemies

	DailyMissionGoldEnemySet9Pos1 = 58 // 25 golden enemies
	DailyMissionGoldEnemySet9Pos2 = 59 // 30 golden enemies
	DailyMissionGoldEnemySet9Pos3 = 60 // 35 golden enemies

	DailyMissionGoldEnemySet10Pos1 = 61 // 25 golden enemies
	DailyMissionGoldEnemySet10Pos2 = 62 // 30 golden enemies
	DailyMissionGoldEnemySet10Pos3 = 63 // 35 golden enemies

	DailyMissionGoldEnemySet11Pos1 = 64 // 30 golden enemies
	DailyMissionGoldEnemySet11Pos2 = 65 // 35 golden enemies
	DailyMissionGoldEnemySet11Pos3 = 66 // 40 golden enemies

	// "Run for x meters!" daily missions
	DailyMissionDistanceSet1Pos1 = 67 // 500 meters
	DailyMissionDistanceSet1Pos2 = 68 // 500 meters
	DailyMissionDistanceSet1Pos3 = 69 // 500 meters

	DailyMissionDistanceSet2Pos1 = 70 // 1000 meters
	DailyMissionDistanceSet2Pos2 = 71 // 1500 meters
	DailyMissionDistanceSet2Pos3 = 72 // 2000 meters

	DailyMissionDistanceSet3Pos1 = 73 // 1500 meters
	DailyMissionDistanceSet3Pos2 = 74 // 2000 meters
	DailyMissionDistanceSet3Pos3 = 75 // 2500 meters

	DailyMissionDistanceSet4Pos1 = 76 // 2500 meters
	DailyMissionDistanceSet4Pos2 = 77 // ? meters
	DailyMissionDistanceSet4Pos3 = 78 // ? meters

	DailyMissionDistanceSet5Pos1 = 79 // 3000 meters
	DailyMissionDistanceSet5Pos2 = 80 // ? meters
	DailyMissionDistanceSet5Pos3 = 81 // ? meters

	DailyMissionDistanceSet6Pos1 = 82 // ? meters
	DailyMissionDistanceSet6Pos2 = 83 // ? meters
	DailyMissionDistanceSet6Pos3 = 84 // ? meters

	DailyMissionDistanceSet7Pos1 = 85 // ? meters
	DailyMissionDistanceSet7Pos2 = 86 // ? meters
	DailyMissionDistanceSet7Pos3 = 87 // ? meters

	DailyMissionDistanceSet8Pos1 = 88 // ? meters
	DailyMissionDistanceSet8Pos2 = 89 // ? meters
	DailyMissionDistanceSet8Pos3 = 90 // ? meters

	DailyMissionDistanceSet9Pos1 = 91 // ? meters
	DailyMissionDistanceSet9Pos2 = 92 // ? meters
	DailyMissionDistanceSet9Pos3 = 93 // ? meters

	DailyMissionDistanceSet10Pos1 = 94 // ? meters
	DailyMissionDistanceSet10Pos2 = 95 // ? meters
	DailyMissionDistanceSet10Pos3 = 96 // ? meters

	DailyMissionDistanceSet11Pos1 = 97 // ? meters
	DailyMissionDistanceSet11Pos2 = 98 // ? meters
	DailyMissionDistanceSet11Pos3 = 99 // ? meters

	// "Save x Animals!" daily missions
	DailyMissionAnimalSet1Pos1 = 100 // 20 animals
	DailyMissionAnimalSet1Pos2 = 101 // 20 animals
	DailyMissionAnimalSet1Pos3 = 102 // 20 animals

	// "Score x points!" daily missions
	DailyMissionScoreSet1Pos1 = 133 // 50000 points
	DailyMissionScoreSet1Pos2 = 134 // 50000 points
	DailyMissionScoreSet1Pos3 = 135 // 50000 points

	// "Collect x Rings!" daily missions
	DailyMissionRingSet1Pos1 = 166 // 300 rings
	DailyMissionRingSet1Pos2 = 167 // 300 rings
	DailyMissionRingSet1Pos3 = 168 // 300 rings
)
