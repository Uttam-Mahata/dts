"""Configuration management for the distributed scheduling system."""

from typing import Dict, Any
import json


class Config:
    """Configuration manager for the scheduling system."""
    
    DEFAULT_CONFIG = {
        'optimizer': {
            'alpha': 0.5,  # Load balancing weight
            'beta': 0.3,   # Resource efficiency weight
            'gamma': 0.2,  # Network latency weight
        },
        'scheduler': {
            'rebalance_threshold': 0.8,
            'heartbeat_timeout': 60.0,
        },
        'system': {
            'max_tasks_per_node': 100,
            'scheduling_interval': 5.0,
        }
    }
    
    def __init__(self, config: Dict[str, Any] = None):
        """
        Initialize configuration.
        
        Args:
            config: Configuration dictionary (uses defaults if None)
        """
        self.config = self.DEFAULT_CONFIG.copy()
        if config:
            self._update_config(self.config, config)
    
    def _update_config(self, base: Dict, updates: Dict):
        """Recursively update configuration."""
        for key, value in updates.items():
            if key in base and isinstance(base[key], dict) and isinstance(value, dict):
                self._update_config(base[key], value)
            else:
                base[key] = value
    
    def get(self, key_path: str, default=None) -> Any:
        """
        Get configuration value by dot-separated path.
        
        Args:
            key_path: Dot-separated path (e.g., 'optimizer.alpha')
            default: Default value if key not found
            
        Returns:
            Configuration value
        """
        keys = key_path.split('.')
        value = self.config
        
        for key in keys:
            if isinstance(value, dict) and key in value:
                value = value[key]
            else:
                return default
        
        return value
    
    def set(self, key_path: str, value: Any):
        """
        Set configuration value by dot-separated path.
        
        Args:
            key_path: Dot-separated path (e.g., 'optimizer.alpha')
            value: Value to set
        """
        keys = key_path.split('.')
        config = self.config
        
        for key in keys[:-1]:
            if key not in config:
                config[key] = {}
            config = config[key]
        
        config[keys[-1]] = value
    
    def to_dict(self) -> Dict[str, Any]:
        """Get configuration as dictionary."""
        return self.config.copy()
    
    def to_json(self) -> str:
        """Get configuration as JSON string."""
        return json.dumps(self.config, indent=2)
    
    @classmethod
    def from_json(cls, json_str: str) -> 'Config':
        """
        Create configuration from JSON string.
        
        Args:
            json_str: JSON string
            
        Returns:
            Config instance
        """
        config_dict = json.loads(json_str)
        return cls(config_dict)
