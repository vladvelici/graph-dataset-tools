function [ similarity ] = sim2( q,z,a,b )
%SIM2 Compute similarity between nodes a and b.
norma = z(a,:)*q*z(a,:)';
normb = z(b,:)*q*z(b,:)';
similarity = norma + normb - 2 * (z(a,:)*q*z(b,:)');
end

