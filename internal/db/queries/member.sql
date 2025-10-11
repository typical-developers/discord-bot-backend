-- name: CreateMemberProfile :one
INSERT INTO guild_profiles (guild_id, member_id)
    VALUES (@guild_id, @member_id)
RETURNING *;

-- name: GetMemberProfile :one
WITH profiles AS (
    SELECT
        member_id,
        card_style,

        -- chat activity info
        ROW_NUMBER() OVER (ORDER BY chat_activity DESC, last_chat_activity_grant DESC) AS chat_activity_rank,
        chat_activity,
        last_chat_activity_grant,

        -- voice activity info
        ROW_NUMBER() OVER (ORDER BY voice_activity DESC, last_voice_activity_grant DESC) AS voice_activity_rank,
        voice_activity,
        last_voice_activity_grant
    FROM guild_profiles
    WHERE
        guild_profiles.guild_id = @guild_id
)
SELECT *
FROM profiles
WHERE
    member_id = @member_id;

-- name: IncrememberMemberChatActivityPoints :one
UPDATE guild_profiles
SET
    chat_activity = chat_activity + @points,
    last_chat_activity_grant = EXTRACT(EPOCH FROM now() AT TIME ZONE 'utc')
WHERE
    guild_id = @guild_id
    AND member_id = @member_id
RETURNING *;

-- name: GetMemberChatActivityRoleInfo :one
WITH
    activity_roles AS (
        SELECT
            role_id,
            SUM(required_points) OVER (
                ORDER BY required_points ASC ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
            ) AS required_points
        FROM guild_activity_roles
        WHERE guild_activity_roles.guild_id = @guild_id
    ),
    all_role_ids AS (
        SELECT CAST(ARRAY_AGG(role_id) AS TEXT[]) AS role_ids
        FROM activity_roles
        WHERE required_points <= CAST(@points AS INT)
    ),
    current_role_info AS (
        SELECT
            role_id,
            required_points
        FROM activity_roles
        WHERE activity_roles.required_points <= CAST(@points AS INT)
        ORDER BY required_points DESC
        LIMIT 1
    ),
    next_role_info AS (
        SELECT
            role_id,
            required_points
        FROM activity_roles
        WHERE activity_roles.required_points > CAST(@points AS INT)
        ORDER BY required_points ASC
        LIMIT 1
    )
SELECT
    all_role_ids.role_ids AS current_roles_ids,
    current_role_info.role_id AS current_role_id,
    current_role_info.required_points AS current_role_required_points,
    next_role_info.role_id AS next_role_id,
    next_role_info.required_points AS next_role_required_points
FROM current_role_info
FULL OUTER JOIN next_role_info ON TRUE
CROSS JOIN all_role_ids;