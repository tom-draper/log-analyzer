type Data = {
    extraction: Extraction[];
    config: Config;
    dataTypes: DataTypes
};

type Extraction = {
    params: Params;
    line: string
    lineNumber: number
    pattern: string
}

type Config = {
    tokens: string[];
    patterns: string[];
    dependencies?: string[];
    conversions?: {
        [token: string]: {
            token: string,
            multiplier: number
        }
    };
}

type DataTypes = {
    [token: string]: {
        [type: string]: number
    }
}

type Params = {
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