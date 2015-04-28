function [ q,z ] = similarity( adj, mu, m )
%SIMILARITY Compute similarity for the undirected graph given.
%   Returns the matrix omega, which can be used to compute the
%   similarities between nodes.
%
%   Arguments:
%   adj     - adjacency matrix
%   miu     - the penalising factor
%   m       - the number of eigenvalues/vectors to use

neigh = sum(adj,2);
neighinv = neigh.^-1;
w = diag(neighinv);
wHalf = diag(sqrt(neighinv));

A = wHalf * adj * wHalf;
[vec, val] = eigs(A,[],m);

disp('Eigenvectors computed.');

gamma = zeros(m,m);
for i=1:m
    gamma(i,i) = (1-mu*val(i,i))^-1;
end

z = diag(sqrt(neigh)) * vec * gamma;
q = vec' * w * vec;

% L = chol(q');
% disp('Cholesky decomposition done.');
% omega = L'*z';

end

