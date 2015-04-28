function [ q,z ] = train( csv_path, mu, k, output_path )
%TRAINTEXT Makes the eigen computation and saves the resulting
% matrices in plain text in the following format:
% - first line contains two integer values: the size of Q
% - the values of Q, space separated
% - two integer values, the size of Z
% - the value of Z

    % some debug info
    fprintf(1, 'csv path: %s\n mu: %s\n k: %s\n output: %s\n',csv_path, mu, k, output_path)

    mu = str2double(mu);
    k = str2double(k);
   
    raw = csvread(csv_path);
    adj = sparse(raw(:,1), raw(:,2), ones(size(raw,1),1));
    [q, z] = similarity(adj, mu, int64(k));

    of = fopen(output_path, 'w');
    fprintf(of, '%d %d\n', size(q,1), size(q,2));
    fprintf(of, '%f ', q);
    fprintf(of, '\n%d %d\n', size(z,1), size(z,2));
    fprintf(of, '%f ', z);
    fclose(of);

    fprintf(1, 'Q(2,3) = %f\n', q(2,3)) 

    % similarity for nodes 3 and 4, for comparison
    fprintf(1, 'Sim 3,4 : %f\n', sim2(q,z,3,4))
end

