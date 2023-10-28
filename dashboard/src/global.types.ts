type Data = {
    extraction: {
        params: LineParams[];
        failed: {
            [lineNumber: string]: string
        };
    };
    config: {
        tokens: string[];
        patterns: string[];
        dependencies?: string[];
        conversions?: {
            [token: string]: {
                token: string,
                multiplier: number
            }
        };
    };
};

type LineParams = {
    [token: string]: string | number;
};

type TokenValueFreq = {
    [token: string]: ValueCount;
};

type ValueCount = {
    [value: string]: number
}

type ValueCounts = {
    [value: string]: number[]
}