CREATE  OR REPLACE FUNCTION getTasksReady(task_status int)
RETURNS TABLE (
    id int,
    doc_id uuid,
    status SMALLINT
)
AS $$
#variable_conflict use_column
BEGIN
    if (task_status IS NULL) then 
        RETURN QUERY
        SELECT dq.id, dq.doc_id, dq.status
        FROM document_queues as dq;
    else
        RETURN QUERY
        SELECT dq.id, dq.doc_id, dq.status
        FROM document_queues as dq
        WHERE status = task_status;
    end if;
END;
$$ LANGUAGE plpgsql;
