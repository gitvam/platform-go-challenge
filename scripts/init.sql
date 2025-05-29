-- Drop existing tables if they exist (for dev resets only)
DROP TABLE IF EXISTS favorites;
DROP TABLE IF EXISTS charts;
DROP TABLE IF EXISTS insights;
DROP TABLE IF EXISTS audiences;

-- Charts table
CREATE TABLE charts (
    id SERIAL PRIMARY KEY,
    external_id TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    x_axis_title TEXT NOT NULL,
    y_axis_title TEXT NOT NULL,
    data INT[] NOT NULL,
    description TEXT NOT NULL
);

-- Insights table
CREATE TABLE insights (
    id SERIAL PRIMARY KEY,
    external_id TEXT NOT NULL UNIQUE,
    text TEXT NOT NULL,
    description TEXT NOT NULL
);

-- Audiences table
CREATE TABLE audiences (
    id SERIAL PRIMARY KEY,
    external_id TEXT NOT NULL UNIQUE,
    gender TEXT NOT NULL,
    birth_country TEXT NOT NULL,
    age_groups TEXT[] NOT NULL,
    hours_on_social INT NOT NULL,
    purchases_last_month INT NOT NULL,
    description TEXT NOT NULL
);

-- Favorites table
CREATE TABLE favorites (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    asset_id INT NOT NULL,
    asset_type TEXT NOT NULL CHECK (asset_type IN ('chart', 'insight', 'audience')),
    description TEXT NOT NULL,
    UNIQUE (user_id, asset_type, asset_id)
);

-- Dummy Charts
INSERT INTO charts (external_id, title, x_axis_title, y_axis_title, data, description) VALUES
  ('chart_engagement_2024', 'Q1 2024 Social Media Engagement', 'Month', 'Engagement (k)', ARRAY[85,92,110,130], 'Tracks monthly engagement for all channels in Q1 2024.'),
  ('chart_ecom_conversion', 'E-commerce Conversion Rates 2024', 'Week', 'Conversion Rate (%)', ARRAY[2,2,3,4,3,5,4], 'Weekly conversion rate trend for Q2 2024.');

-- Dummy Insights
INSERT INTO insights (external_id, text, description) VALUES
  ('insight_active_users', '78% of millennials engage with branded content daily.', 'Based on 2024 survey data across EMEA.'),
  ('insight_genz_tiktok', 'Gen Z users are 3x more likely to purchase after seeing a TikTok ad.', 'Finding from global digital consumer study 2024.');

-- Dummy Audiences
INSERT INTO audiences (external_id, gender, birth_country, age_groups, hours_on_social, purchases_last_month, description) VALUES
  ('aud_greece_men_24_35', 'male', 'Greece', ARRAY['24-35'], 4, 3, 'Digitally active Greek men aged 24-35 with high purchasing intent.'),
  ('aud_uk_females_18_24', 'female', 'UK', ARRAY['18-24'], 6, 5, 'UK-based young women, highly active on Instagram and TikTok.');

-- Dummy favorites (cross-ref using IDs)
INSERT INTO favorites (user_id, asset_id, asset_type, description) VALUES
  ('11111111-1111-1111-1111-111111111111', 1, 'chart', 'Tracks monthly engagement for all channels in Q1 2024.'),
  ('11111111-1111-1111-1111-111111111111', 1, 'insight', 'Based on 2024 survey data across EMEA.'),
  ('11111111-1111-1111-1111-111111111111', 1, 'audience', 'Digitally active Greek men aged 24-35 with high purchasing intent.'),
  ('22222222-2222-2222-2222-222222222222', 2, 'chart', 'Weekly conversion rate trend for Q2 2024.'),
  ('22222222-2222-2222-2222-222222222222', 2, 'insight', 'Finding from global digital consumer study 2024.'),
  ('22222222-2222-2222-2222-222222222222', 2, 'audience', 'UK-based young women, highly active on Instagram and TikTok.');