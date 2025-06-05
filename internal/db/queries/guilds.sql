-- name: GetAllTimeChatActivityRankings :many
SELECT
    CAST (
        ROW_NUMBER() OVER (ORDER BY guild_profiles.activity_points DESC) AS INT
    ) AS rank,
    member_id,
    activity_points
FROM guild_profiles
WHERE
    guild_id = @guild_id
ORDER BY activity_points DESC
LIMIT 15
OFFSET @offset_by;