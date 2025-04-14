// API endpoints
const API_BASE_URL = 'http://localhost:8080/api';
const REFRESH_INTERVAL = 5000; // 5 seconds

// Function to fetch top players
async function fetchLeaderboard() {
    try {
        const response = await fetch(`${API_BASE_URL}/leaderboard/top`);
        if (!response.ok) throw new Error('Failed to fetch leaderboard');
        
        const data = await response.json();
        updateLeaderboardTable(data.leaderboard);
    } catch (error) {
        console.error('Error fetching leaderboard:', error);
    }
}

// Function to look up a player's rank
async function lookupPlayer() {
    const playerId = document.getElementById('player-id').value;
    if (!playerId) {
        showPlayerResult('Please enter a player ID', false);
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/leaderboard/rank/${playerId}`);
        if (!response.ok) throw new Error('Player not found');
        
        const data = await response.json();
        showPlayerResult(`
            Player: ${data.username}
            Rank: #${data.rank}
            Total Score: ${data.total_score}
        `, true);
    } catch (error) {
        showPlayerResult('Player not found or an error occurred', false);
    }
}

// Function to update the leaderboard table
function updateLeaderboardTable(leaderboard) {
    const tbody = document.getElementById('leaderboard-body');
    tbody.innerHTML = '';

    leaderboard.forEach(player => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>#${player.rank}</td>
            <td>${player.username}</td>
            <td>${player.total_score}</td>
        `;
        tbody.appendChild(row);
    });
}

// Function to display player lookup results
function showPlayerResult(message, isSuccess) {
    const resultDiv = document.getElementById('player-result');
    resultDiv.textContent = message;
    resultDiv.className = `player-result ${isSuccess ? 'success' : 'error'}`;
}

// Start periodic updates
fetchLeaderboard();
setInterval(fetchLeaderboard, REFRESH_INTERVAL);