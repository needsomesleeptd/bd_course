

CREATE OR REPLACE FUNCTION addToQueue()
Returns TRIGGER
AS $$
BEGIN
    IF NOT EXISTS (
        SELECT
            *
        FROM
            document_queues
        WHERE
            doc_id = NEW.ID AND status = 2 
    ) AND NEW.document_name NOT LIKE '%VIP%' 
    THEN
        INSERT INTO 
        document_queues(doc_id,status)
        VALUES
        (NEW.ID,0); -- status zero means unchecked
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insertQueueTrigger AFTER INSERT ON public.documents
FOR EACH ROW EXECUTE PROCEDURE addToQueue();



CREATE  OR REPLACE FUNCTION getTasksReady(task_status SMALLINT)
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

select * from getTasksReady();

INSERT INTO documents (id, page_count, document_name, checks_count, creator_id, creation_time) 
VALUES 
('710eb9fc-1118-4a66-837e-8e3d2d027a68', 52, 'NIRS.pdf', 0, 5, '2024-05-10 13:48:59.541053');