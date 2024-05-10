CREATE  OR REPLACE FUNCTION getTasksReady()
RETURNS TABLE (
    id int,
    doc_id uuid,
    status SMALLINT
)
AS $$
#variable_conflict use_column
BEGIN
    RETURN QUERY
    SELECT dq.id, dq.doc_id, dq.status
    FROM document_queues as dq
    WHERE status = 0;
END;
$$ LANGUAGE plpgsql;

select * from getTasksReady();
