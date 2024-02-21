DROP TABLE IF EXISTS public.clientes;
CREATE TABLE public.clientes (
  id SERIAL PRIMARY KEY NOT NULL,
  nome varchar(100) NOT NULL,
  limite INT NOT NULL,
  saldo INT NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS public.transacoes;
CREATE TABLE public.transacoes (
  id SERIAL PRIMARY KEY NOT NULL,
  tipo char(1) NOT NULL,
  valor INT NOT NULL,
  saldo_pos_transacao INT,
  descricao varchar(10),
  realizado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  cliente_id INT,
  CONSTRAINT fk_cliente
    FOREIGN KEY(cliente_id)
      REFERENCES clientes(id)
);


INSERT INTO clientes (nome, limite)
  VALUES
    ('o barato sai caro', 1000 * 100),
    ('zan corp ltda', 800 * 100),
    ('les cruders', 10000 * 100),
    ('padaria joia de cocaia', 100000 * 100),
    ('kid mais', 5000 * 100);