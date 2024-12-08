CREATE EXTENSION pg_stat_statements;

load 'auto_explain';
SET auto_explain.log_min_duration = 400;
SET auto_explain.log_analyze = true;
