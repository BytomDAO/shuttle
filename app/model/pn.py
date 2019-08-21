# Appendix: improved modulo calc

def inv(n, q):
    """div on PN modulo a/b mod q as a * inv(b, q) mod q
    >>> assert n * inv(n, q) % q == 1
    """
    # n*inv % q = 1 => n*inv = q*m + 1 => n*inv + q*-m = 1
    # => egcd(n, q) = (inv, -m, 1) => inv = egcd(n, q)[0] (mod q)
    return egcd(n, q)[0] % q
    #[ref] naive implementation
    #for i in range(q):
    #    if (n * i) % q == 1:
    #        return i
    #    pass
    #assert False, "unreached"
    #pass


def egcd(a, b):
    """extended GCD
    returns: (s, t, gcd) as a*s + b*t == gcd
    >>> s, t, gcd = egcd(a, b)
    >>> assert a % gcd == 0 and b % gcd == 0
    >>> assert a * s + b * t == gcd
    """
    s0, s1, t0, t1 = 1, 0, 0, 1
    while b > 0:
        q, r = divmod(a, b)
        a, b = b, r
        s0, s1, t0, t1 = s1, s0 - q * s1, t1, t0 - q * t1
        pass
    return s0, t0, a

def inv2(n, q):
    """another PN invmod: from euler totient function
    - n ** (q - 1) % q = 1 => n ** (q - 2) % q = n ** -1 % q
    """
    assert q > 2
    s, p2, p = 1, n, q - 2
    while p > 0:
        if p & 1 == 1: s = s * p2 % q
        p, p2 = p >> 1, pow(p2, 2, q)
        pass
    return s


def sqrt(n, q):
    """sqrt on PN modulo: returns two numbers or exception if not exist
    >>> assert (sqrt(n, q)[0] ** 2) % q == n
    >>> assert (sqrt(n, q)[1] ** 2) % q == n
    """
    assert n < q
    for i in range(1, q):
        if pow(i, 2, q) == n:
            return (i, q - i)
        pass
    raise Exception("not found")


def sqrt2(n, q):
    """sqrtmod for bigint
    - Algorithm 3.34 of http://www.cacr.math.uwaterloo.ca/hac/about/chap3.pdf
    """
    import random
    # b: some non-quadratic-residue
    b = 0 
    while b == 0 or jacobi(b, q) != -1:
        b = random.randint(1, q - 1)
        pass
    # q = t * 2^s + 1, t is odd
    t, s = q - 1, 0 
    while t & 1 == 0:
        t, s = t >> 1, s + 1
        pass
    assert q == t * pow(2, s) + 1 and t % 2 == 1
    ni = inv(n, q)
    c = pow(b, t, q)
    r = pow(n, (t + 1) // 2, q)
    for i in range(1, s):
        d = pow(pow(r, 2, q) * ni % q, pow(2, s - i - 1, q), q)
        if d == q - 1: r = r * c % q
        c = pow(c, 2, q)
        pass
    return (r, q - r)


def jacobi(a, q):
    """jacobi symbol: judge existing sqrtmod (1: exist, -1: not exist)
    - j(a*b,q) = j(a,q)*j(b,q)
    - j(a*q+b, q) = j(b, q)
    - j(a, 1) = 1
    - j(0, q) = 0
    - j(2, q) = -1 ** (q^2 - 1)/8
    - j(p, q) = -1 ^ {(p - 1)/2 * (q - 1)/2} * j(q, p)
    """
    if q == 1: return 1
    if a == 0: return 0
    if a % 2 == 0: return (-1) ** ((q * q - 1) // 8) * jacobi(a // 2, q)
    return (-1) ** ((a - 1) // 2 * (q - 1) // 2) * jacobi(q % a, a)

def jacobi2(a, q):
    """quick jacobi symbol
    - algorithm 2.149 of http://www.cacr.math.uwaterloo.ca/hac/about/chap2.pdf
    """
    if a == 0: return 0
    if a == 1: return 1
    a1, e = a, 0
    while a1 & 1 == 0:
        a1, e = a1 >> 1, e + 1
        pass
    m8 = q % 8
    s = -1 if m8 == 3 or m8 == 5 else 1 # m8 = 0,2,4,6 and 1,7
    if q % 4 == 3 and a1 % 4 == 3: s = -s
    return s if a1 == 1 else s * jacobi2(q % a1, a1)

    