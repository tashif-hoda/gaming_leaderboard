// API endpoints
const API_BASE_URL = 'http://localhost:8080/api';
const REFRESH_INTERVAL = 5000; // 5 seconds

let retryTimeout;

// Function to fetch top players
async function fetchLeaderboard() {
    try {
        const response = await fetch(`${API_BASE_URL}/leaderboard/top`);
        if (response.status === 429) {
            handleRateLimit(response);
            return;
        }
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
        if (response.status === 429) {
            handleRateLimit(response);
            return;
        }
        if (!response.ok) throw new Error('Player not found');
        
        const data = await response.json();
        showPlayerResult(`
            Rank: #${data.rank}
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

// Function to handle rate limit responses
async function handleRateLimit(response) {
    const retryAfter = response.headers.get('Retry-After') || 5;
    const message = `Rate limit exceeded. Please wait ${retryAfter} seconds before trying again.`;
    showPlayerResult(message, false);
    
    // Clear any existing retry timeout
    if (retryTimeout) {
        clearTimeout(retryTimeout);
    }
    
    // Set retry timeout
    retryTimeout = setTimeout(() => {
        fetchLeaderboard();
    }, retryAfter * 1000);
}

// Start periodic updates
fetchLeaderboard();
setInterval(fetchLeaderboard, REFRESH_INTERVAL);